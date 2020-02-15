package params

import (
	"time"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/verification"

	libparams "github.com/qdm12/golibs/params"
)

type ParamsGetter interface {
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

type paramsGetter struct {
	envParams   libparams.EnvParams
	verifier    verification.Verifier
	fileManager files.FileManager
}

func NewParamsGetter(envParams libparams.EnvParams) ParamsGetter {
	return &paramsGetter{
		envParams:   envParams,
		verifier:    verification.NewVerifier(),
		fileManager: files.NewFileManager(),
	}
}

// GetOutputDir obtains the output directory path to write files to
// from the environment variable OUTPUT_DIR and defaults to ./files
func (p *paramsGetter) GetOutputDir() (path string, err error) {
	return p.envParams.GetPath("OUTPUT_DIR", libparams.Default("./files"))
}

// GetPeriod obtains the period in minutes from the PERIOD environment
// variable. It defaults to 600 minutes.
func (p *paramsGetter) GetPeriod() (periodMinutes time.Duration, err error) {
	return p.envParams.GetDuration("PERIOD", libparams.Default("600m"))
}
