package dnscrypto

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/qdm12/golibs/network"
)

// GetNamedRoot downloads the named.root and returns it
func GetNamedRoot(httpClient *http.Client, namedRootHexMD5 string) (namedRoot []byte, err error) {
	namedRoot, err = network.GetContent(
		httpClient,
		"https://www.internic.net/domain/named.root",
		network.GetContentParamsType{DisguisedUserAgent: true})
	if err != nil {
		return nil, err
	}
	sum := md5.Sum(namedRoot)
	hexSum := hex.EncodeToString(sum[:])
	if hexSum != namedRootHexMD5 {
		return nil, fmt.Errorf("named root MD5 sum %q is not expected sum %q", hexSum, namedRootHexMD5)
	}
	return namedRoot, nil
}
