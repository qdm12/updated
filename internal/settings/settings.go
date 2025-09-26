// Package settings handles the application settings.
package settings

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
	"github.com/qdm12/updated/pkg/dnscrypto"
)

// Settings holds the application settings.
type Settings struct {
	OutputDir        string
	Period           time.Duration
	ResolveHostnames *bool
	HTTPTimeout      time.Duration
	HexSums          struct {
		NamedRootMD5      *string
		RootAnchorsSHA256 string
	}
	Git      Git
	Log      Log
	Shoutrrr Shoutrrr
}

func (s *Settings) Read(r *reader.Reader) (err error) {
	s.OutputDir = r.String("OUTPUT_DIR")

	s.Period, err = r.Duration("PERIOD")
	if err != nil {
		return err
	}

	s.ResolveHostnames, err = r.BoolPtr("RESOLVE_HOSTNAMES")
	if err != nil {
		return err
	}

	s.HTTPTimeout, err = r.Duration("HTTP_TIMEOUT")
	if err != nil {
		return err
	}

	s.HexSums.NamedRootMD5 = r.Get("NAMED_ROOT_MD5")
	s.HexSums.RootAnchorsSHA256 = r.String("ROOT_ANCHORS_SHA256")

	err = s.Git.read(r)
	if err != nil {
		return fmt.Errorf("reading git settings: %w", err)
	}

	s.Log.Read(r)
	s.Shoutrrr.read(r)

	return nil
}

var (
	regex32BytesHex = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	regex64BytesHex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)

	ErrPeriodTooSmall            = errors.New("period is too small")
	ErrNamedRootMD5NotValid      = errors.New("named root MD5 checksum is not valid")
	ErrRootAnchorsSHA256NotValid = errors.New("root anchors SHA256 checksum is not valid")
)

// SetDefaults sets the default values for the settings.
func (s *Settings) SetDefaults() {
	s.OutputDir = gosettings.DefaultComparable(s.OutputDir, "./files")
	const defaultPeriod = 600 * time.Minute
	s.Period = gosettings.DefaultComparable(s.Period, defaultPeriod)
	s.ResolveHostnames = gosettings.DefaultPointer(s.ResolveHostnames, false)
	const defaultHTTPTimeout = 10 * time.Second
	s.HTTPTimeout = gosettings.DefaultComparable(s.HTTPTimeout, defaultHTTPTimeout)
	s.HexSums.NamedRootMD5 = gosettings.DefaultPointer(s.HexSums.NamedRootMD5, "")
	s.HexSums.RootAnchorsSHA256 = gosettings.DefaultComparable(s.HexSums.RootAnchorsSHA256, dnscrypto.RootAnchorsSHA256Sum)
	s.Git.setDefaults()
	s.Log.SetDefaults()
	s.Shoutrrr.setDefaults()
}

// Validate validates the settings and returns an error if something is wrong.
func (s Settings) Validate() (err error) {
	_, err = filepath.Abs(s.OutputDir)
	if err != nil {
		return fmt.Errorf("output directory: %w", err)
	}

	const minPeriod = 5 * time.Minute
	switch {
	case s.Period < minPeriod:
		return fmt.Errorf("%w: %s < %s", ErrPeriodTooSmall, s.Period, minPeriod)
	case *s.HexSums.NamedRootMD5 != "" &&
		!regex64BytesHex.MatchString(*s.HexSums.NamedRootMD5):
		return fmt.Errorf("%w: %q does not match regex %q",
			ErrNamedRootMD5NotValid, *s.HexSums.NamedRootMD5, regex32BytesHex)
	case !regex64BytesHex.MatchString(s.HexSums.RootAnchorsSHA256):
		return fmt.Errorf("%w: %q does not match regex %q",
			ErrRootAnchorsSHA256NotValid, s.HexSums.RootAnchorsSHA256, regex64BytesHex)
	}

	const minHTTPTimeout = time.Second
	if s.HTTPTimeout < minHTTPTimeout {
		return fmt.Errorf("HTTP timeout %s cannot be smaller than %s",
			s.HTTPTimeout, minHTTPTimeout)
	}

	err = s.Git.validate()
	if err != nil {
		return fmt.Errorf("validating git settings: %w", err)
	}

	err = s.Log.Validate()
	if err != nil {
		return fmt.Errorf("validating log settings: %w", err)
	}

	err = s.Shoutrrr.validate()
	if err != nil {
		return fmt.Errorf("validating shoutrrr settings: %w", err)
	}

	return nil
}

func (s Settings) String() string {
	return s.toLinesNode().String()
}

func (s Settings) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Settings summary:")
	node.Appendf("output directory: %s", s.OutputDir)
	node.Appendf("period: %s", s.Period)
	node.Appendf("resolve hostnames: %s", gosettings.BoolToYesNo(s.ResolveHostnames))
	node.Appendf("HTTP timeout: %s", s.HTTPTimeout)
	node.Appendf("named root MD5 sum: %s", *s.HexSums.NamedRootMD5)
	node.Appendf("root anchors SHA256 sum: %s", s.HexSums.RootAnchorsSHA256)
	node.AppendNode(s.Git.toLinesNode())
	node.AppendNode(s.Log.toLinesNode())
	node.AppendNode(s.Shoutrrr.toLinesNode())
	return node
}
