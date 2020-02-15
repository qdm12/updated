package settings

import (
	"fmt"
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

func Get(paramsGetter params.ParamsGetter) (s Settings, err error) {
	s.OutputDir, err = paramsGetter.GetOutputDir()
	if err != nil {
		return s, err
	}
	s.HexSums.NamedRootMD5, err = paramsGetter.GetNamedRootMD5()
	if err != nil {
		return s, err
	}
	s.HexSums.RootAnchorsSHA256, err = paramsGetter.GetRootAnchorsSHA256()
	if err != nil {
		return s, err
	}
	s.Period, err = paramsGetter.GetPeriod()
	if err != nil {
		return s, err
	}
	s.ResolveHostnames, err = paramsGetter.GetResolveHostnames()
	if err != nil {
		return s, err
	}
	git, err := paramsGetter.GetGit()
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
	s.Git.GitURL, err = paramsGetter.GetGitURL()
	if err != nil {
		return s, err
	}
	s.Git.SSHKnownHosts, err = paramsGetter.GetSSHKnownHostsFilepath()
	if err != nil {
		return s, err
	}
	s.Git.SSHKey, err = paramsGetter.GetSSHKeyFilepath()
	if err != nil {
		return s, err
	}
	s.Git.SSHKeyPassphrase, err = paramsGetter.GetSSHKeyPassphrase()
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
	kv := map[string]string{
		"output directory":        s.OutputDir,
		"period":                  s.Period.String(),
		"resolve hostnames":       resolveHostnamesStr,
		"named root MD5 sum":      s.HexSums.NamedRootMD5,
		"root anchors SHA256 sum": s.HexSums.RootAnchorsSHA256,
	}
	if s.Git == nil {
		kv["Git"] = "disabled"
	} else {
		kv["Git URL"] = s.Git.GitURL
		kv["SSH known hosts file"] = s.Git.SSHKnownHosts
		kv["SSH key file"] = s.Git.SSHKey
		kv["SSH key passphrase file"] = s.Git.SSHKeyPassphrase
	}
	result = "Settings:\n"
	for k, v := range kv {
		result += fmt.Sprintf("|-- %s: %s\n", k, v)
	}
	return result
}
