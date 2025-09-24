// Package dnscrypto provides functionality for downloading and verifying
// DNS cryptographic files such as the named root and root anchors.
package dnscrypto

import (
	"net/http"
	"sync"
)

// DNSCrypto handles downloading and verifying DNS crypto related files.
type DNSCrypto struct {
	client                 *http.Client
	namedRootHexMD5        string
	namedRootHexMD5Mu      sync.RWMutex
	rootAnchorsHexSHA256   string
	rootAnchorsHexSHA256Mu sync.RWMutex
}

// New creates a new DNSCrypto object.
func New(client *http.Client,
	namedRootHexMD5, rootAnchorsHexSHA256 string,
) *DNSCrypto {
	return &DNSCrypto{
		client:               client,
		namedRootHexMD5:      namedRootHexMD5,
		rootAnchorsHexSHA256: rootAnchorsHexSHA256,
	}
}
