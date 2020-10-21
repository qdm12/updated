package params

import (
	"fmt"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/updated/pkg/constants"
)

// GetNamedRootMD5 obtains the MD5 Hex encoded checksum for the named root
// from the environment variable NAMED_ROOT_MD5.
func (p *getter) GetNamedRootMD5() (namedRootMD5 string, err error) {
	s, err := p.envParams.GetEnv("NAMED_ROOT_MD5", libparams.Default(constants.NamedRootMD5Sum))
	if err != nil {
		return "", err
	} else if !p.verifier.MatchMD5String(s) {
		return "", fmt.Errorf("%s is not a 32 hexadecimal MD5 string", s)
	}
	return s, nil
}

// GetRootAnchorsSHA256 obtains the SHA256 Hex encoded checksum for the root anchors
// from the environment variable ROOT_ANCHORS_SHA256.
func (p *getter) GetRootAnchorsSHA256() (rootAnchorsSHA256 string, err error) {
	s, err := p.envParams.GetEnv("ROOT_ANCHORS_SHA256", libparams.Default(constants.RootAnchorsSHA256Sum))
	if err != nil {
		return "", err
	} else if !p.verifier.Match64BytesHex(s) {
		return "", fmt.Errorf("%s is not a 64 hexadecimal SHA256 string", s)
	}
	return s, nil
}
