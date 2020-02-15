package params

import (
	"fmt"

	"github.com/qdm12/updated/pkg/constants"

	libparams "github.com/qdm12/golibs/params"
)

// GetNamedRootMD5 obtains the MD5 Hex encoded checksum for the named root
// from the environment variable NAMED_ROOT_MD5. It defaults to
// 25659425b11bb58ece6306d0cfe4b587
func (p *paramsGetter) GetNamedRootMD5() (namedRootMD5 string, err error) {
	s, err := p.envParams.GetEnv("NAMED_ROOT_MD5", libparams.Default(constants.NamedRootMD5Sum))
	if err != nil {
		return "", err
	} else if !p.verifier.MatchMD5String(s) {
		return "", fmt.Errorf("%s is not a 32 hexadecimal MD5 string", s)
	}
	return s, nil
}

// GetRootAnchorsSHA256 obtains the SHA256 Hex encoded checksum for the root anchors
// from the environment variable ROOT_ANCHORS_SHA256. It defaults to
// 45336725f9126db810a59896ae93819de743c416262f79c4444042c92e520770
func (p *paramsGetter) GetRootAnchorsSHA256() (rootAnchorsSHA256 string, err error) {
	s, err := p.envParams.GetEnv("ROOT_ANCHORS_SHA256", libparams.Default(constants.RootAnchorsSHA256Sum))
	if err != nil {
		return "", err
	} else if !p.verifier.Match64BytesHex(s) {
		return "", fmt.Errorf("%s is not a 64 hexadecimal SHA256 string", s)
	}
	return s, nil
}
