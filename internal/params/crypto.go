package params

import (
	"errors"
	"fmt"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/updated/pkg/dnscrypto"
)

var ErrNamedRootMD5Invalid = errors.New("named root MD5 sum is invalid")

// GetNamedRootMD5 obtains the MD5 Hex encoded checksum for the named root
// from the environment variable NAMED_ROOT_MD5.
func (p *getter) GetNamedRootMD5() (namedRootMD5 string, err error) {
	s, err := p.envParams.Get("NAMED_ROOT_MD5")
	switch {
	case err != nil:
		return "", err
	case s == "":
		return "", nil
	case !p.verifier.MatchMD5String(s):
		return "", fmt.Errorf("%w: not a 32 hexadecimal string: %s", ErrNamedRootMD5Invalid, s)
	}
	return s, nil
}

var ErrRootAnchorsSHA256Invalid = errors.New("root anchors SHA256 sum is invalid")

// GetRootAnchorsSHA256 obtains the SHA256 Hex encoded checksum for the root anchors
// from the environment variable ROOT_ANCHORS_SHA256.
func (p *getter) GetRootAnchorsSHA256() (rootAnchorsSHA256 string, err error) {
	s, err := p.envParams.Get("ROOT_ANCHORS_SHA256", libparams.Default(dnscrypto.RootAnchorsSHA256Sum))
	if err != nil {
		return "", err
	} else if !p.verifier.Match64BytesHex(s) {
		return "", fmt.Errorf("%w: not a 64 hexadecimal string: %s", ErrRootAnchorsSHA256Invalid, s)
	}
	return s, nil
}
