package dnscrypto

//nolint:gci
import (
	"context"
	"crypto/md5" //nolint:gosec
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/qdm12/updated/pkg/constants"
)

// GetNamedRoot downloads the named.root and returns it.
func (d *dnsCrypto) GetNamedRoot(ctx context.Context) (namedRoot []byte, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, constants.NamedRootURL, nil)
	if err != nil {
		return nil, err
	}
	response, err := d.client.Do(request)
	if err != nil {
		return nil, err
	} else if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("%w: %d %s", ErrBadStatusCode, response.StatusCode, response.Status)
	}

	namedRoot, err = ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	sum := md5.Sum(namedRoot) //nolint:gosec
	hexSum := hex.EncodeToString(sum[:])
	if hexSum != d.namedRootHexMD5 {
		return nil, fmt.Errorf("%w: %q is not expected %q", ErrChecksumMismatch, hexSum, d.namedRootHexMD5)
	}

	return namedRoot, nil
}
