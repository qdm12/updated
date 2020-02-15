package params

import (
	"fmt"
	"regexp"

	libparams "github.com/qdm12/golibs/params"
)

// GetGit obtains 'yes' or 'no' to do Git operations, from the environment
// variable GIT, and defaults to no.
func (p *paramsGetter) GetGit() (doGit bool, err error) {
	return p.envParams.GetYesNo("GIT", libparams.Default("no"))
}

// GetSSHKnownHostsFilepath obtains the file path of the SSH known_hosts file,
// from the environment variable SSH_KNOWN_HOSTS and defaults to /known_hosts.
func (p *paramsGetter) GetSSHKnownHostsFilepath() (filePath string, err error) {
	filePath, err = p.envParams.GetPath("SSH_KNOWN_HOSTS", libparams.Default("./known_hosts"))
	if err != nil {
		return "", err
	}
	exists, err := p.fileManager.FileExists(filePath)
	if err != nil {
		return "", err
	} else if !exists {
		return "", fmt.Errorf("SSH known hosts file %q does not exist", filePath)
	}
	return filePath, nil
}

// GetSSHKeyFilepath obtains the file path of the SSH private key,
// from the environment variable SSH_KEY and defaults to /key
func (p *paramsGetter) GetSSHKeyFilepath() (filePath string, err error) {
	filePath, err = p.envParams.GetPath("SSH_KEY", libparams.Default("./key"))
	if err != nil {
		return "", err
	}
	exists, err := p.fileManager.FileExists(filePath)
	if err != nil {
		return "", err
	} else if !exists {
		return "", fmt.Errorf("SSH key file %q does not exist", filePath)
	}
	return filePath, nil
}

// GetSSHKeyPassphrase obtains the SSH key passphrase file path,
// from the environment variable SSH_KEY_PASSPHRASE and defaults to returning an
// empty string passphrase if no file is provided.
// It uses files instead of environment variables for security reasons.
func (p *paramsGetter) GetSSHKeyPassphrase() (passphrase string, err error) {
	filePath, err := p.envParams.GetPath("SSH_KEY_PASSPHRASE")
	if err != nil {
		return "", err
	}
	if filePath == "" {
		// no passphrase
		return "", nil
	}
	exists, err := p.fileManager.FileExists(filePath)
	if err != nil {
		return "", err
	} else if !exists {
		return "", fmt.Errorf("SSH passphrase file %q does not exist", filePath)
	}
	data, err := p.fileManager.ReadFile(filePath)
	return string(data), err
}

// GetGitURL obtains the Git repository URL to interact with,
// from the environment variable GIT_URL.
func (p *paramsGetter) GetGitURL() (URL string, err error) {
	u, err := p.envParams.GetURL("GIT_URL")
	URL = u.String()
	if !regexp.MustCompile(`((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`).MatchString(URL) {
		return "", fmt.Errorf("environment variable GIT_URL value %q is not valid", URL)
	}
	return URL, nil
}
