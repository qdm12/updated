package dnscrypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/golibs/network"
)

// RootAnchors holds the XML data of the root anchors
type RootAnchors struct {
	TrustAnchor struct {
		ID        string `xml:"id,attr"`
		Source    string `xml:"source,attr"`
		Zone      string `xml:"Zone"`
		KeyDigest []struct {
			ID         string    `xml:"id,attr"`
			ValidFrom  time.Time `xml:"validFrom,attr"`
			ValidUntil time.Time `xml:"validUntil,attr"`
			KeyTag     int       `xml:"KeyTag"`
			Algorithm  int       `xml:"Algorithm"`
			DigestType int       `xml:"DigestType"`
			Digest     string    `xml:"Digest"`
		} `xml:"KeyDigest"`
	} `xml:"TrustAnchor"`
}

// GetRootAnchorsXML fetches the root anchors XML file online and parses it
func GetRootAnchorsXML(httpClient *http.Client, rootAnchorsHexSHA256 string) (rootAnchorsXML []byte, err error) {
	rootAnchorsXML, err = network.GetContent(
		httpClient,
		"https://data.iana.org/root-anchors/root-anchors.xml",
		network.GetContentParamsType{DisguisedUserAgent: true})
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(rootAnchorsXML)
	hexSum := hex.EncodeToString(sum[:])
	if hexSum != rootAnchorsHexSHA256 {
		return nil, fmt.Errorf("root anchors SHA256 sum %q is not expected sum %q", hexSum, rootAnchorsHexSHA256)
	}
	return rootAnchorsXML, err
}

// ConvertRootAnchorsToRootKeys converts root anchors XML data to a list of DNS root keys
func ConvertRootAnchorsToRootKeys(rootAnchorsXML []byte) (rootKeys []string, err error) {
	rootAnchors, err := parseRootAnchors(rootAnchorsXML)
	if err != nil {
		return nil, err
	}
	for _, keyDigest := range rootAnchors.TrustAnchor.KeyDigest {
		rootKey := fmt.Sprintf(". IN DS %d %d %d %s",
			keyDigest.KeyTag, keyDigest.Algorithm, keyDigest.DigestType, keyDigest.Digest)
		rootKeys = append(rootKeys, rootKey)
	}
	return rootKeys, nil
}

func parseRootAnchors(content []byte) (rootAnchorsXML RootAnchors, err error) {
	err = xml.Unmarshal(content, &rootAnchorsXML)
	if err != nil {
		return rootAnchorsXML, err
	}
	return rootAnchorsXML, nil
}
