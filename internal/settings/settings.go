// Package settings handles the application settings.
package settings

import (
	"fmt"
	"strings"
	"time"
)

// Settings holds the application settings.
type Settings struct {
	OutputDir        string
	Period           time.Duration
	ResolveHostnames bool
	HexSums          struct {
		NamedRootMD5      string
		RootAnchorsSHA256 string
	}
	Git *Git
}

// Git holds the Git related settings.
type Git struct {
	GitURL           string
	SSHKnownHosts    string
	SSHKey           string
	SSHKeyPassphrase string
}

// Getter defines an interface to get settings.
type Getter interface {
	// General getters
	GetOutputDir() (path string, err error)
	GetPeriod() (period time.Duration, err error)

	// Git
	GetGit() (doGit bool, err error)
	GetSSHKnownHostsFilepath() (filePath string, err error)
	GetSSHKeyFilepath() (filePath string, err error)
	GetSSHKeyPassphrase() (passphrase string, err error)
	GetGitURL() (URL string, err error)

	// Crypto
	GetNamedRootMD5() (namedRootMD5 string, err error)
	GetRootAnchorsSHA256() (rootAnchorsSHA256 string, err error)

	// IPs blocking
	GetResolveHostnames() (resolveHostnames bool, err error)
}

// Get retrieves the settings using the provided [Getter].
func Get(getter Getter) (s Settings, err error) {
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
	} else if git {
		s.Git = new(Git)
		*s.Git, err = readGit(getter)
		if err != nil {
			return s, fmt.Errorf("reading git settings: %w", err)
		}
	}

	return s, nil
}

func readGit(getter Getter) (g Git, err error) {
	g.GitURL, err = getter.GetGitURL()
	if err != nil {
		return g, fmt.Errorf("getting Git URL: %w", err)
	}
	g.SSHKnownHosts, err = getter.GetSSHKnownHostsFilepath()
	if err != nil {
		return g, fmt.Errorf("getting SSH known hosts filepath: %w", err)
	}
	g.SSHKey, err = getter.GetSSHKeyFilepath()
	if err != nil {
		return g, fmt.Errorf("getting SSH key filepath: %w", err)
	}
	g.SSHKeyPassphrase, err = getter.GetSSHKeyPassphrase()
	if err != nil {
		return g, fmt.Errorf("getting SSH key passphrase: %w", err)
	}
	return g, nil
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
