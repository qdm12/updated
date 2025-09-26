package settings

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// Git holds the Git related settings.
type Git struct {
	Enabled          *bool
	GitURL           string
	SSHKnownHosts    string
	SSHKey           string
	SSHKeyPassphrase *string
}

func (l *Git) read(r *reader.Reader) (err error) {
	l.Enabled, err = r.BoolPtr("GIT")
	if err != nil {
		return err
	}

	l.GitURL = r.String("GIT_URL")
	l.SSHKnownHosts = r.String("SSH_KNOWN_HOSTS")
	l.SSHKey = r.String("SSH_KEY")
	l.SSHKeyPassphrase = r.Get("SSH_KEY_PASSPHRASE")

	return nil
}

func (l *Git) setDefaults() {
	l.Enabled = gosettings.DefaultPointer(l.Enabled, true)
	l.SSHKnownHosts = gosettings.DefaultComparable(l.SSHKnownHosts, "./known_hosts")
	l.SSHKey = gosettings.DefaultComparable(l.SSHKey, "./key")
	l.SSHKeyPassphrase = gosettings.DefaultPointer(l.SSHKeyPassphrase, "")
}

var (
	regexpGitURL = regexp.MustCompile(`((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`)

	ErrGitURLNotValid = errors.New("git URL is not valid")
)

func (l Git) validate() (err error) {
	if !*l.Enabled {
		return nil
	}

	switch {
	case l.GitURL == "":
		return fmt.Errorf("%w: empty", ErrGitURLNotValid)
	case !regexpGitURL.MatchString(l.GitURL):
		return fmt.Errorf("%w: %s", ErrGitURLNotValid, l.GitURL)
	}

	err = checkFileExists(l.SSHKnownHosts)
	if err != nil {
		return fmt.Errorf("checking SSH known hosts file: %w", err)
	}

	err = checkFileExists(l.SSHKey)
	if err != nil {
		return fmt.Errorf("checking SSH key file: %w", err)
	}

	if *l.SSHKeyPassphrase != "" {
		err = checkFileExists(*l.SSHKeyPassphrase)
		if err != nil {
			return fmt.Errorf("checking SSH key passphrase file: %w", err)
		}
	}

	return nil
}

func (l Git) toLinesNode() (node *gotree.Node) {
	if !*l.Enabled {
		return gotree.New("Git: disabled")
	}
	node = gotree.New("Git settings:")
	node.Appendf("Git URL: %s", l.GitURL)
	node.Appendf("SSH known hosts file: %s", l.SSHKnownHosts)
	node.Appendf("SSH key file: %s", l.SSHKey)
	if *l.SSHKeyPassphrase == "" {
		node.Appendf("SSH key passphrase: [not set]")
	} else {
		node.Appendf("SSH key passphrase: %s", *l.SSHKeyPassphrase)
	}
	return node
}

func checkFileExists(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0) //nolint:gosec
	if os.IsNotExist(err) {
		return fmt.Errorf("%w: filepath %q", err, filePath)
	} else if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("closing file %q: %w", filePath, err)
	}
	return nil
}
