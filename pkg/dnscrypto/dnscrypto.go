package dnscrypto

import (
	"context"
	"net/http"
)

type DNSCrypto interface {
	DownloadNamedRoot(ctx context.Context) (namedRoot []byte, err error)
	DownloadRootAnchorsXML(ctx context.Context) (rootAnchorsXML []byte, err error)
	ConvertRootAnchorsToRootKeys(rootAnchorsXML []byte) (rootKeys []string, err error)
}

type dnsCrypto struct {
	client               *http.Client
	namedRootHexMD5      string
	rootAnchorsHexSHA256 string
}

func New(client *http.Client, namedRootHexMD5, rootAnchorsHexSHA256 string) DNSCrypto {
	return &dnsCrypto{
		client:               client,
		namedRootHexMD5:      namedRootHexMD5,
		rootAnchorsHexSHA256: rootAnchorsHexSHA256}
}
