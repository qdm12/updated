package dnscrypto

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/qdm12/golibs/network"
	"github.com/qdm12/updated/pkg/constants"
)

// GetNamedRoot downloads the named.root and returns it
func (d *dnsCrypto) GetNamedRoot() (namedRoot []byte, err error) {
	namedRoot, status, err := d.client.GetContent(constants.NamedRootURL, network.UseRandomUserAgent())
	if err != nil {
		return nil, err
	} else if status != http.StatusOK {
		return nil, fmt.Errorf("HTTP status is %d", status)
	}
	sum := md5.Sum(namedRoot)
	hexSum := hex.EncodeToString(sum[:])
	if hexSum != d.namedRootHexMD5 {
		return nil, fmt.Errorf("named root MD5 sum %q is not expected sum %q", hexSum, d.namedRootHexMD5)
	}
	return namedRoot, nil
}
