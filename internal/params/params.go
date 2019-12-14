package params

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/qdm12/golibs/verification"

	libparams "github.com/qdm12/golibs/params"
)

// GetOutputDir obtains the output directory path to write files to
// from the environment variable OUTPUT_DIR and defaults to ./files
func GetOutputDir() (path string, err error) {
	return libparams.GetPath("OUTPUT_DIR", "./files")
}

// GetNamedRootMD5 obtains the MD5 Hex encoded checksum for the named root
// from the environment variable NAMED_ROOT_MD5. It defaults to
// 25659425b11bb58ece6306d0cfe4b587
func GetNamedRootMD5() (namedRootMD5 string, err error) {
	s := libparams.GetEnv("NAMED_ROOT_MD5", "1e4e7c3e1ce2c5442eed998046edf548")
	if !verification.MatchMD5String(s) {
		return "", fmt.Errorf("%s is not a 32 hexadecimal MD5 string", s)
	}
	return s, nil
}

// GetRootAnchorsSHA256 obtains the SHA256 Hex encoded checksum for the root anchors
// from the environment variable ROOT_ANCHORS_SHA256. It defaults to
// 45336725f9126db810a59896ae93819de743c416262f79c4444042c92e520770
func GetRootAnchorsSHA256() (rootAnchorsSHA256 string, err error) {
	s := libparams.GetEnv("ROOT_ANCHORS_SHA256", "45336725f9126db810a59896ae93819de743c416262f79c4444042c92e520770")
	if !verification.MatchSHA256String(s) {
		return "", fmt.Errorf("%s is not a 64 hexadecimal SHA256 string", s)
	}
	return s, nil
}

// GetPeriodMinutes obtains the period in minutes from the PERIOD environment
// variable. It defaults to 600 minutes.
func GetPeriodMinutes() (periodMinutes time.Duration, err error) {
	duration, err := libparams.GetDuration("PERIOD", 600, time.Minute)
	if err != nil {
		return periodMinutes, err
	}
	return duration, nil
}

// GetResolveHostnames obtains 'yes' or 'no' to resolve hostnames in order to obtain
// more IP addresses, from the environment variable RESOLVE_HOSTNAMES, and defaults to no.
// If you are blocking the hostname resolution on your network, turn this feature off.
func GetResolveHostnames() (resolveHostnames bool, err error) {
	s := libparams.GetEnv("RESOLVE_HOSTNAMES", "no")
	if s == "yes" {
		return true, nil
	} else if s == "no" {
		return false, nil
	}
	return false, fmt.Errorf("RESOLVE_HOSTNAMES value %q can only be 'yes' or 'no'", s)
}

// GetGit obtains 'yes' or 'no' to do Git operations, from the environment
// variable GIT, and defaults to no.
func GetGit() (doGit bool, err error) {
	s := libparams.GetEnv("GIT", "no")
	if s == "yes" {
		return true, nil
	} else if s == "no" {
		return false, nil
	}
	return false, fmt.Errorf("GIT value %q can only be 'yes' or 'no'", s)
}

// GetSSHKnownHostsFilepath obtains the file path of the SSH known_hosts file,
// from the environment variable SSH_KNOWN_HOSTS and defaults to /known_hosts.
func GetSSHKnownHostsFilepath() (filePath string, err error) {
	filePath, err = libparams.GetPath("SSH_KNOWN_HOSTS", "./known_hosts")
	if err != nil {
		return "", err
	} else if _, err := os.Stat(filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

// GetSSHKeyFilepath obtains the file path of the SSH private key,
// from the environment variable SSH_KEY and defaults to /key
func GetSSHKeyFilepath() (filePath string, err error) {
	filePath, err = libparams.GetPath("SSH_KEY", "./key")
	if err != nil {
		return "", err
	} else if _, err := os.Stat(filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

// GetSSHKeyPassphrase obtains the SSH key passphrase file path,
// from the environment variable SSH_KEY_PASSPHRASE and defaults to returning an
// empty string passphrase if no file is provided.
// It uses files instead of environment variables for security reasons.
func GetSSHKeyPassphrase() (passphrase string, err error) {
	filePath, err := libparams.GetPath("SSH_KEY_PASSPHRASE", "")
	if err != nil {
		return "", err
	}
	if filePath == "" {
		// no passphrase
		return "", nil
	}
	if _, err := os.Stat(filePath); err != nil {
		return "", err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetGitURL obtains the Git repository URL to interact with,
// from the environment variable GIT_URL.
func GetGitURL() (URL string, err error) {
	URL = libparams.GetEnv("GIT_URL", "")
	if !regexp.MustCompile(`((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`).MatchString(URL) {
		return "", fmt.Errorf("environment variable GIT_URL value %q is not valid", URL)
	}
	return URL, nil
}
