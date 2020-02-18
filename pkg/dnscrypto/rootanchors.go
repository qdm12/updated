package dnscrypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/golibs/network"
	"github.com/qdm12/updated/pkg/constants"
)

// TrustAnchor holds the XML data of the root anchors
type TrustAnchor struct {
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

// GetRootAnchorsXML fetches the root anchors XML file online and parses it
func (d *dnsCrypto) GetRootAnchorsXML() (rootAnchorsXML []byte, err error) {
	rootAnchorsXML, status, err := d.client.GetContent(constants.RootAnchorsURL, network.UseRandomUserAgent())
	if err != nil {
		return nil, err
	} else if status != http.StatusOK {
		return nil, fmt.Errorf("HTTP status is %d", status)
	}
	sum := sha256.Sum256(rootAnchorsXML)
	hexSum := hex.EncodeToString(sum[:])
	if hexSum != d.rootAnchorsHexSHA256 {
		return nil, fmt.Errorf("root anchors SHA256 sum %q is not expected sum %q", hexSum, d.rootAnchorsHexSHA256)
	}
	return rootAnchorsXML, err
}

// ConvertRootAnchorsToRootKeys converts root anchors XML data to a list of DNS root keys
func (d *dnsCrypto) ConvertRootAnchorsToRootKeys(rootAnchorsXML []byte) (rootKeys []string, err error) {
	var trustAnchor TrustAnchor
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
