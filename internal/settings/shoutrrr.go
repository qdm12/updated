package settings

import (
	"errors"
	"fmt"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// Shoutrrr holds the shoutrrr related settings.
type Shoutrrr struct {
	ServiceURLs []string
}

func (s *Shoutrrr) read(r *reader.Reader) {
	s.ServiceURLs = r.CSV("SHOUTRRR_SERVICES", reader.ForceLowercase(false))
}

func (s *Shoutrrr) setDefaults() {}

var ErrShoutrrrServiceURLNotValid = errors.New("shoutrrr service URL is not valid")

func (s Shoutrrr) validate() (err error) {
	router := &router.ServiceRouter{}
	for _, url := range s.ServiceURLs {
		_, _, err := router.ExtractServiceName(url)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrShoutrrrServiceURLNotValid, err)
		}
	}
	return nil
}

func (s Shoutrrr) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Shoutrrr services:")
	for _, url := range s.ServiceURLs {
		node.Appendf("%s", url)
	}
	return node
}
