package settings

import (
	"errors"
	"fmt"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// Log holds the logging related settings.
type Log struct {
	Level string
}

func (s *Log) Read(r *reader.Reader) {
	s.Level = r.String("LOG_LEVEL")
}

// SetDefaults sets the default values for the log settings.
func (s *Log) SetDefaults() {
	s.Level = gosettings.DefaultComparable(s.Level, logging.LevelInfo.String())
}

var ErrLogLevelNotValid = errors.New("log level is not valid")

// Validate validates the log settings.
func (s Log) Validate() (err error) {
	switch s.Level {
	case logging.LevelDebug.String(), logging.LevelInfo.String(),
		logging.LevelWarn.String(), logging.LevelError.String():
		return nil
	default:
		return fmt.Errorf("%w: %q", ErrLogLevelNotValid, s.Level)
	}
}

func (s Log) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Log:")
	node.Appendf("Level: %s", s.Level)
	return node
}
