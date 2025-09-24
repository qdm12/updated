// Package params handles obtaining parameters from the environment.
package params

import (
	"time"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
)

// Getter holds methods to obtain parameters from the environment.
type Getter struct {
	envParams libparams.Interface
	verifier  verification.Verifier
}

// NewGetter creates a new Getter.
func NewGetter(envParams libparams.Interface) *Getter {
	return &Getter{
		envParams: envParams,
		verifier:  verification.NewVerifier(),
	}
}

// GetOutputDir obtains the output directory path to write files to
// from the environment variable OUTPUT_DIR and defaults to ./files.
func (p *Getter) GetOutputDir() (path string, err error) {
	return p.envParams.Path("OUTPUT_DIR", libparams.Default("./files"))
}

// GetPeriod obtains the period in minutes from the PERIOD environment
// variable. It defaults to 600 minutes.
func (p *Getter) GetPeriod() (periodMinutes time.Duration, err error) {
	return p.envParams.Duration("PERIOD", libparams.Default("600m"))
}
