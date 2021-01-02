package dnscrypto

import (
	"context"
	"net/http"
	"sync"
)

type DNSCrypto interface {
	DownloadNamedRoot(ctx context.Context) (namedRoot []byte, err error)
	DownloadRootAnchorsXML(ctx context.Context) (rootAnchorsXML []byte, err error)
	ConvertRootAnchorsToRootKeys(rootAnchorsXML []byte) (rootKeys []string, err error)
	SetNamedRootHexMD5(namedRootHexMD5 string)
	SetRootAnchorsHexSHA256(rootAnchorsHexSHA256 string)
}

type dnsCrypto struct {
	client                 *http.Client
	namedRootHexMD5Mu      sync.RWMutex
	rootAnchorsHexSHA256Mu sync.RWMutex
	namedRootHexMD5        string
	rootAnchorsHexSHA256   string
}

func New(client *http.Client, namedRootHexMD5, rootAnchorsHexSHA256 string) DNSCrypto {
	return &dnsCrypto{
		client:               client,
		namedRootHexMD5:      namedRootHexMD5,
		rootAnchorsHexSHA256: rootAnchorsHexSHA256}
}
