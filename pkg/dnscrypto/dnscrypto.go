package dnscrypto

import (
	"context"
	"net/http"
	"sync"
)

var _ Interface = (*DNSCrypto)(nil)

type Interface interface {
	NamedRootManager
	RootAnchorsManager
}

type NamedRootManager interface {
	SetNamedRootHexMD5(namedRootHexMD5 string)
	DownloadNamedRoot(ctx context.Context) (namedRoot []byte, err error)
}
type RootAnchorsManager interface {
	SetRootAnchorsHexSHA256(rootAnchorsHexSHA256 string)
	DownloadRootAnchorsXML(ctx context.Context) (rootAnchorsXML []byte, err error)
	ConvertRootAnchorsToRootKeys(rootAnchorsXML []byte) (rootKeys []string, err error)
}

type DNSCrypto struct {
	client                 *http.Client
	namedRootHexMD5        string
	namedRootHexMD5Mu      sync.RWMutex
	rootAnchorsHexSHA256   string
	rootAnchorsHexSHA256Mu sync.RWMutex
}

func New(client *http.Client,
	namedRootHexMD5, rootAnchorsHexSHA256 string) *DNSCrypto {
	return &DNSCrypto{
		client:               client,
		namedRootHexMD5:      namedRootHexMD5,
		rootAnchorsHexSHA256: rootAnchorsHexSHA256}
}
