package dnscrypto

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DownloadRootAnchorsXML fetches the root anchors XML file online and parses it.
func (d *DNSCrypto) DownloadRootAnchorsXML(ctx context.Context) (rootAnchorsXML []byte, err error) {
	const url = "https://data.iana.org/root-anchors/root-anchors.xml"
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	rootAnchorsXML, err = io.ReadAll(response.Body)
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

// ConvertRootAnchorsToRootKeys converts root anchors XML data
// to a list of DNS root keys.
func (d *DNSCrypto) ConvertRootAnchorsToRootKeys(rootAnchorsXML []byte) (
	rootKeys []string, err error) {
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
	err = xml.Unmarshal(rootAnchorsXML, &trustAnchor)
	if err != nil {
		return nil, err
	}

	rootKeys = make([]string, len(trustAnchor.KeyDigest))
	for i, keyDigest := range trustAnchor.KeyDigest {
		rootKeys[i] = fmt.Sprintf(". IN DS %d %d %d %s", keyDigest.KeyTag,
			keyDigest.Algorithm, keyDigest.DigestType, keyDigest.Digest)
	}

	return rootKeys, nil
}
