package params

import (
	"time"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
)

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

type getter struct {
	envParams libparams.Interface
	verifier  verification.Verifier
}

func NewGetter(envParams libparams.Interface) Getter {
	return &getter{
		envParams: envParams,
		verifier:  verification.NewVerifier(),
	}
}

// GetOutputDir obtains the output directory path to write files to
// from the environment variable OUTPUT_DIR and defaults to ./files.
func (p *getter) GetOutputDir() (path string, err error) {
	return p.envParams.Path("OUTPUT_DIR", libparams.Default("./files"))
}

// GetPeriod obtains the period in minutes from the PERIOD environment
// variable. It defaults to 600 minutes.
func (p *getter) GetPeriod() (periodMinutes time.Duration, err error) {
	return p.envParams.Duration("PERIOD", libparams.Default("600m"))
}
