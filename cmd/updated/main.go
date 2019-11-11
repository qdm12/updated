package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kyokomi/emoji"

	"github.com/qdm12/updated/internal/env"
	"github.com/qdm12/updated/internal/params"
	"github.com/qdm12/updated/pkg/dnscrypto"
	"github.com/qdm12/updated/pkg/git"
	"github.com/qdm12/updated/pkg/hostnames"
	"github.com/qdm12/updated/pkg/ips"

	"github.com/qdm12/golibs/admin"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/healthcheck"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/signals"
)

func main() {
	if healthcheck.Mode(os.Args) {
		if err := healthcheck.Query(); err != nil {
			logging.Err(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	fmt.Println("#####################################")
	fmt.Println("########## Updated #########")
	fmt.Println("########## by Quentin McGaw #########")
	fmt.Println("########## Give some " + emoji.Sprint(":heart:") + "at ##########")
	fmt.Println("# github.com/qdm12/updated #")
	fmt.Print("#####################################\n\n")
	encoding, level, nodeID, err := libparams.GetLoggerConfig()
	if err != nil {
		logging.Error(err.Error())
	} else {
		logging.InitLogger(encoding, level, nodeID)
	}
	var e env.Env
	HTTPTimeout, err := libparams.GetHTTPTimeout(3000)
	e.HTTPClient = &http.Client{Timeout: HTTPTimeout}
	e.Gotify = admin.InitGotify(e.HTTPClient)
	outputDir, err := params.GetOutputDir()
	e.FatalOnError(err)
	namedRootMD5, err := params.GetNamedRootMD5()
	e.FatalOnError(err)
	rootAnchorsSHA256, err := params.GetRootAnchorsSHA256()
	e.FatalOnError(err)
	periodMinutes, err := params.GetPeriodMinutes()
	e.FatalOnError(err)
	resolveHostnames, err := params.GetResolveHostnames()
	e.FatalOnError(err)
	doGit, err := params.GetGit()
	e.FatalOnError(err)
	var knownHostsPath, keyPath, keyPassphrase, gitURL string
	if doGit {
		knownHostsPath, err = params.GetSSHKnownHostsFilepath()
		e.FatalOnError(err)
		keyPath, err = params.GetSSHKeyFilepath()
		e.FatalOnError(err)
		keyPassphrase, err = params.GetSSHKeyPassphrase()
		e.FatalOnError(err)
		gitURL, err = params.GetGitURL()
		e.FatalOnError(err)
	}
	go signals.WaitForExit(e.ShutdownFromSignal)
	errs := network.ConnectivityChecks(e.HTTPClient, []string{"google.com"})
	for _, err := range errs {
		e.Warn(err)
	}
	e.Gotify.Notify("Program started", 1, "")
	for {
		go run(
			e.HTTPClient,
			e.CheckError,
			outputDir,
			periodMinutes,
			namedRootMD5,
			rootAnchorsSHA256,
			resolveHostnames,
			doGit,
			knownHostsPath,
			keyPath,
			keyPassphrase,
			gitURL,
		)
		time.Sleep(periodMinutes)
	}
}

func run(httpClient *http.Client, checkOnError func(err error) error, outputDir string, periodMinutes time.Duration,
	namedRootMD5, rootAnchorsSHA256 string, resolveHostnames, doGit bool, knownHostsPath, keyPath, keyPassphrase, gitURL string) {
	tStart := time.Now()
	defer func() {
		executionTime := time.Since(tStart)
		logging.Infof("overall execution took %s", executionTime)
		logging.Infof("sleeping for %s", periodMinutes-executionTime)
	}()
	if doGit {
		// Setup Git repository
		gitClient, err := git.NewClient(knownHostsPath, keyPath, keyPassphrase, gitURL, outputDir)
		if checkOnError(err) != nil {
			return
		}
		err = gitClient.CheckoutBranch("updated")
		if err != nil {
			err := gitClient.Branch("updated")
			if checkOnError(err) != nil {
				return
			}
		}
		err = gitClient.CheckoutBranch("updated")
		if checkOnError(err) != nil {
			return
		}
		err = gitClient.Pull()
		if checkOnError(err) != nil {
			return
		}
		// Commit changes and upload on branch updated
		defer func() {
			err := gitClient.Pull()
			if checkOnError(err) != nil {
				return
			}
			// TODO
			// err := gitClient.Commit()
			// err := gitClient.Push()
		}()
	}

	// Build named root from internic.net
	namedRoot, err := dnscrypto.GetNamedRoot(httpClient, namedRootMD5)
	if checkOnError(err) != nil {
		err = files.WriteToFile(filepath.Join(outputDir, "named.root.updated"), namedRoot)
		checkOnError(err)
	}

	// Build root anchors XML from data.iana.org
	rootAnchorsXML, err := dnscrypto.GetRootAnchorsXML(httpClient, rootAnchorsSHA256)
	if checkOnError(err) != nil {
		err = files.WriteToFile(filepath.Join(outputDir, "root-anchors.xml.updated"), rootAnchorsXML)
		checkOnError(err)
	}
	rootKeys, err := dnscrypto.ConvertRootAnchorsToRootKeys(rootAnchorsXML)
	if checkOnError(err) != nil {
		err = files.WriteLinesToFile(filepath.Join(outputDir, "root-anchors.xml.updated"), rootKeys)
		checkOnError(err)
	}

	// Build hostnames
	var hostnamesFailed bool
	hostnames, err := hostnames.BuildMalicious(httpClient)
	if checkOnError(err) != nil {
		hostnamesFailed = true
		err = files.WriteLinesToFile(filepath.Join(outputDir, "malicious-hostnames.updated"), hostnames)
		checkOnError(err)
	}

	IPs := []string{}
	if !hostnamesFailed && resolveHostnames {
		IPs = append(IPs, ips.BuildIPsFromHostnames(hostnames)...)
	}

	// Build IPs
	newIPs, err := ips.BuildMalicious(httpClient)
	if checkOnError(err) != nil {
		IPs = append(IPs, newIPs...)
		var removedCount int
		IPs, removedCount = ips.CleanIPs(IPs)
		logging.Infof("Trimmed down %d IP address lines", removedCount)
		err = files.WriteLinesToFile(filepath.Join(outputDir, "malicious-ips.updated"), IPs)
		checkOnError(err)
	}
}
