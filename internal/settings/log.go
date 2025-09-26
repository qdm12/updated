package settings

import (
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
	"github.com/qdm12/log"
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
	s.Level = gosettings.DefaultComparable(s.Level, log.LevelInfo.String())
}

// Validate validates the log settings.
func (s Log) Validate() (err error) {
	_, err = log.ParseLevel(s.Level)
	return err
}

func (s Log) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Log:")
	node.Appendf("Level: %s", s.Level)
	return node
}
