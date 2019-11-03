package params

import (
	"fmt"
	"time"

	"github.com/qdm12/golibs/verification"

	libparams "github.com/qdm12/golibs/params"
)

// GetOutputDir obtains the output directory path to write files to
// from the environment variable OUTPUT_DIR and defaults to ./files
func GetOutputDir() (path string, err error) {
	return libparams.GetPath("OUTPUT_DIR", "./files")
}

// GetNamedRootMD5 obtains the MD5 Hex encoded checksum for the named root
// from the environment variable NAMED_ROOT_MD5. It defaults to
// 25659425b11bb58ece6306d0cfe4b587
func GetNamedRootMD5() (namedRootMD5 string, err error) {
	s := libparams.GetEnv("NAMED_ROOT_MD5", "23ec4e704cdaa1dcaaa6f66bc2c0563f")
	if !verification.MatchMD5String(s) {
		return "", fmt.Errorf("%s is not a 32 hexadecimal MD5 string", s)
	}
	return s, nil
}

// GetRootAnchorsSHA256 obtains the SHA256 Hex encoded checksum for the root anchors
// from the environment variable ROOT_ANCHORS_SHA256. It defaults to
// 45336725f9126db810a59896ae93819de743c416262f79c4444042c92e520770
func GetRootAnchorsSHA256() (rootAnchorsSHA256 string, err error) {
	s := libparams.GetEnv("ROOT_ANCHORS_SHA256", "45336725f9126db810a59896ae93819de743c416262f79c4444042c92e520770")
	if !verification.MatchSHA256String(s) {
		return "", fmt.Errorf("%s is not a 64 hexadecimal SHA256 string", s)
	}
	return s, nil
}

// GetPeriodMinutes obtains the period in minutes from the PERIOD environment
// variable. It defaults to 600 minutes.
func GetPeriodMinutes() (periodMinutes time.Duration, err error) {
	duration, err := libparams.GetDuration("PERIOD", 600, time.Minute)
	if err != nil {
		return periodMinutes, err
	}
	return duration, nil
}

// GetResolveHostnames obtains 'yes' or 'no' to resolve hostnames in order to obtain
// more IP addresses, from the environment variable RESOLVE_HOSTNAMES.
// If you are blocking the hostname resolution on your network, turn this feature off.
func GetResolveHostnames() (resolveHostnames bool, err error) {
	s := libparams.GetEnv("RESOLVE_HOSTNAMES", "no")
	if s == "yes" {
		return true, nil
	} else if s == "no" {
		return false, nil
	}
	return false, fmt.Errorf("RESOLVE_HOSTNAMES value %q can only be 'yes' or 'no'", s)
}
