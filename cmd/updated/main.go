package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kyokomi/emoji"

	"github.com/qdm12/updated/internal/env"
	"github.com/qdm12/updated/internal/params"
	"github.com/qdm12/updated/pkg/dnscrypto"
	"github.com/qdm12/updated/pkg/hostnames"
	"github.com/qdm12/updated/pkg/ips"

	"github.com/qdm12/golibs/admin"
	"github.com/qdm12/golibs/healthcheck"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/signals"
)

func main() {
	logging.InitLogger()
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
	var e env.Env
	HTTPTimeout, err := libparams.GetHTTPTimeout(3 * time.Second)
	e.HTTPClient = &http.Client{Timeout: HTTPTimeout}
	e.Gotify = admin.InitGotify(e.HTTPClient)
	outputDir, err := params.GetOutputDir()
	e.FatalOnError(err)
	// Create output dir if it does not exist
	// TODO replace with golibs/files
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0700)
		e.FatalOnError(err)
	}
	namedRootMD5, err := params.GetNamedRootMD5()
	e.FatalOnError(err)
	rootAnchorsSHA256, err := params.GetRootAnchorsSHA256()
	e.FatalOnError(err)
	periodMinutes, err := params.GetPeriodMinutes()
	e.FatalOnError(err)
	resolveHostnames, err := params.GetResolveHostnames()
	e.FatalOnError(err)
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
			resolveHostnames)
		time.Sleep(periodMinutes)
	}
}

func run(httpClient *http.Client, checkOnError func(err error), outputDir string, periodMinutes time.Duration,
	namedRootMD5, rootAnchorsSHA256 string, resolveHostnames bool) {
	tStart := time.Now()
	defer func() {
		executionTime := time.Since(tStart)
		logging.Infof("overall execution took %s", executionTime)
		logging.Infof("sleeping for %s", periodMinutes-executionTime)
	}()

	// Build named root from internic.net
	namedRoot, err := dnscrypto.GetNamedRoot(httpClient, namedRootMD5)
	checkOnError(err)
	err = writeToFile(outputDir, "named.root.updated", namedRoot)
	checkOnError(err)

	// Build root anchors XML from data.iana.org
	rootAnchorsXML, err := dnscrypto.GetRootAnchorsXML(httpClient, rootAnchorsSHA256)
	checkOnError(err)
	err = writeToFile(outputDir, "root-anchors.xml.updated", rootAnchorsXML)
	checkOnError(err)
	rootKeys, err := dnscrypto.ConvertRootAnchorsToRootKeys(rootAnchorsXML)
	checkOnError(err)
	err = writeToFile(outputDir, "root-anchors.xml.updated", convertLinesToBytes(rootKeys))
	checkOnError(err)

	// Build hostnames
	hostnames, err := hostnames.BuildMalicious(httpClient)
	checkOnError(err)
	err = writeToFile(outputDir, "malicious-hostnames.updated", convertLinesToBytes(hostnames))
	checkOnError(err)

	// Build IPs from hostnames
	IPs := []string{}
	if resolveHostnames {
		IPs = append(IPs, ips.BuildIPsFromHostnames(hostnames)...)
	}

	// Build IPs
	newIPs, err := ips.BuildMalicious(httpClient)
	checkOnError(err)
	IPs = append(IPs, newIPs...)
	var removedCount int
	IPs, removedCount = ips.CleanIPs(IPs)
	logging.Infof("Trimmed down %d IP address lines", removedCount)
	err = writeToFile(outputDir, "malicious-ips.updated", convertLinesToBytes(IPs))
	checkOnError(err)
}

// TODO replace with golibs/files
func writeToFile(outputDir, outputFilename string, data []byte) error {
	// TODO create dir if not existing
	targetPath := filepath.Clean(outputDir + "/" + outputFilename)
	return ioutil.WriteFile(targetPath, data, 0700)
}

func convertLinesToBytes(lines []string) []byte {
	s := strings.Join(lines, "\n")
	return []byte(s)
}
