package params

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	libparams "github.com/qdm12/golibs/params"
)

var (
	ErrSSHKnownHostFileDoesNotExist = errors.New("SSH known hosts file does not exist")
	ErrSSHKeyFileDoesNotExist       = errors.New("SSH key file does not exist")
)

// GetGit obtains 'yes' or 'no' to do Git operations, from the environment
// variable GIT, and defaults to no.
func (p *getter) GetGit() (doGit bool, err error) {
	return p.envParams.YesNo("GIT", libparams.Default("no"))
}

// GetSSHKnownHostsFilepath obtains the file path of the SSH known_hosts file,
// from the environment variable SSH_KNOWN_HOSTS and defaults to /known_hosts.
func (p *getter) GetSSHKnownHostsFilepath() (filePath string, err error) {
	filePath, err = p.envParams.Path("SSH_KNOWN_HOSTS", libparams.Default("./known_hosts"))
	if err != nil {
		return "", err
	}

	file, err := p.osOpenFile(filePath, os.O_RDONLY, 0)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("%w: filepath %q", ErrSSHKnownHostFileDoesNotExist, filePath)
	} else if err != nil {
		return "", err
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	return filePath, nil
}

// GetSSHKeyFilepath obtains the file path of the SSH private key,
// from the environment variable SSH_KEY and defaults to ./key.
func (p *getter) GetSSHKeyFilepath() (filePath string, err error) {
	filePath, err = p.envParams.Path("SSH_KEY", libparams.Default("./key"))
	if err != nil {
		return "", err
	}

	file, err := p.osOpenFile(filePath, os.O_RDONLY, 0)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("%w: filepath %q", ErrSSHKnownHostFileDoesNotExist, filePath)
	} else if err != nil {
		return "", err
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	return filePath, nil
}

// GetSSHKeyPassphrase obtains the SSH key passphrase file path,
// from the environment variable SSH_KEY_PASSPHRASE and defaults to returning an
// empty string passphrase if no file is provided.
// It uses files instead of environment variables for security reasons.
func (p *getter) GetSSHKeyPassphrase() (passphrase string, err error) {
	filePath, err := p.envParams.Path("SSH_KEY_PASSPHRASE")
	if err != nil {
		return "", err
	}
	if filePath == "" {
		// no passphrase
		return "", nil
	}

	file, err := p.osOpenFile(filePath, os.O_RDONLY, 0)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("%w: filepath %q", ErrSSHKeyFileDoesNotExist, filePath)
	} else if err != nil {
		return "", err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	return string(b), nil
}

var ErrInvalidGitURL = errors.New("invalid git URL")

// GetGitURL obtains the Git repository URL to interact with,
// from the environment variable GIT_URL.
func (p *getter) GetGitURL() (url string, err error) {
	url, err = p.envParams.Get("GIT_URL", libparams.Compulsory())
	if err != nil {
		return "", err
	}
	if !regexp.MustCompile(`((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`).MatchString(url) {
		return "", fmt.Errorf("%w: from environment variable GIT_URL: %s", ErrInvalidGitURL, url)
	}
	return url, nil
}
