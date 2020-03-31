package settings

import (
	"strings"
	"time"

	"github.com/qdm12/updated/internal/params"
)

type Settings struct {
	OutputDir        string
	Period           time.Duration
	ResolveHostnames bool
	HexSums          struct {
		NamedRootMD5      string
		RootAnchorsSHA256 string
	}
	Git *struct {
		GitURL           string
		SSHKnownHosts    string
		SSHKey           string
		SSHKeyPassphrase string
	}
}

func Get(getter params.Getter) (s Settings, err error) {
	s.OutputDir, err = getter.GetOutputDir()
	if err != nil {
		return s, err
	}
	s.HexSums.NamedRootMD5, err = getter.GetNamedRootMD5()
	if err != nil {
		return s, err
	}
	s.HexSums.RootAnchorsSHA256, err = getter.GetRootAnchorsSHA256()
	if err != nil {
		return s, err
	}
	s.Period, err = getter.GetPeriod()
	if err != nil {
		return s, err
	}
	s.ResolveHostnames, err = getter.GetResolveHostnames()
	if err != nil {
		return s, err
	}
	git, err := getter.GetGit()
	if err != nil {
		return s, err
	} else if !git {
		return s, nil
	}
	s.Git = new(struct {
		GitURL           string
		SSHKnownHosts    string
		SSHKey           string
		SSHKeyPassphrase string
	})
	s.Git.GitURL, err = getter.GetGitURL()
	if err != nil {
		return s, err
	}
	s.Git.SSHKnownHosts, err = getter.GetSSHKnownHostsFilepath()
	if err != nil {
		return s, err
	}
	s.Git.SSHKey, err = getter.GetSSHKeyFilepath()
	if err != nil {
		return s, err
	}
	s.Git.SSHKeyPassphrase, err = getter.GetSSHKeyPassphrase()
	if err != nil {
		return s, err
	}
	return s, nil
}

func (s *Settings) String() (result string) {
	resolveHostnamesStr := "no"
	if s.ResolveHostnames {
		resolveHostnamesStr = "yes"
	}
	lines := []string{
		"output directory: " + s.OutputDir,
		"period: " + s.Period.String(),
		"resolve hostnames: " + resolveHostnamesStr,
		"named root MD5 sum: " + s.HexSums.NamedRootMD5,
		"root anchors SHA256 sum: " + s.HexSums.RootAnchorsSHA256,
	}
	if s.Git == nil {
		lines = append(lines, "Git: disabled")
	} else {
		passhpraseSet := "no"
		if len(s.Git.SSHKeyPassphrase) > 0 {
			passhpraseSet = "yes"
		}
		lines = append(lines, []string{
			"Git URL: " + s.Git.GitURL,
			"SSH known hosts file: " + s.Git.SSHKnownHosts,
			"SSH key file: " + s.Git.SSHKey,
			"SSH key passphrase set: " + passhpraseSet,
		}...)
	}
	return "Settings:\n|--" + strings.Join(lines, "\n|--") + "\n"
}
