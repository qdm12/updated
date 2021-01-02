package dnscrypto

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/qdm12/updated/pkg/constants"
)

// DownloadRootAnchorsXML fetches the root anchors XML file online and parses it.
func (d *dnsCrypto) DownloadRootAnchorsXML(ctx context.Context) (rootAnchorsXML []byte, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, constants.RootAnchorsURL, nil)
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
	rootAnchorsXML, err = ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	checksum := d.getRootAnchorsHexSHA256()
	if checksum == "" {
		return rootAnchorsXML, nil
	}

	sum := sha256.Sum256(rootAnchorsXML)
	hexSum := hex.EncodeToString(sum[:])
	if hexSum != checksum {
		return nil, fmt.Errorf("%w: %q is not expected %q", ErrChecksumMismatch, hexSum, checksum)
	}

	return rootAnchorsXML, nil
}

// ConvertRootAnchorsToRootKeys converts root anchors XML data to a list
// of DNS root keys.
func (d *dnsCrypto) ConvertRootAnchorsToRootKeys(rootAnchorsXML []byte) (rootKeys []string, err error) {
	var trustAnchor struct {
		XMLName   xml.Name `xml:"TrustAnchor"`
		ID        string   `xml:"id,attr"`
		Source    string   `xml:"source,attr"`
		Zone      string   `xml:"Zone"`
		KeyDigest []struct {
			ID         string    `xml:"id,attr"`
			ValidFrom  time.Time `xml:"validFrom,attr"`
			ValidUntil time.Time `xml:"validUntil,attr"`
			KeyTag     int       `xml:"KeyTag"`
			Algorithm  int       `xml:"Algorithm"`
			DigestType int       `xml:"DigestType"`
			Digest     string    `xml:"Digest"`
		} `xml:"KeyDigest"`
	}
	if err := xml.Unmarshal(rootAnchorsXML, &trustAnchor); err != nil {
		return nil, err
	}
	for _, keyDigest := range trustAnchor.KeyDigest {
		rootKey := fmt.Sprintf(". IN DS %d %d %d %s",
			keyDigest.KeyTag, keyDigest.Algorithm, keyDigest.DigestType, keyDigest.Digest)
		rootKeys = append(rootKeys, rootKey)
	}
	return rootKeys, nil
}
