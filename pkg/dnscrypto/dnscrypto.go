package dnscrypto

import "github.com/qdm12/golibs/network"

type DNSCrypto interface {
	GetNamedRoot() (namedRoot []byte, err error)
	GetRootAnchorsXML() (rootAnchorsXML []byte, err error)
	ConvertRootAnchorsToRootKeys(rootAnchorsXML []byte) (rootKeys []string, err error)
}

type dnsCrypto struct {
	client               network.Client
	namedRootHexMD5      string
	rootAnchorsHexSHA256 string
}

func NewDNSCrypto(client network.Client, namedRootHexMD5, rootAnchorsHexSHA256 string) DNSCrypto {
	return &dnsCrypto{
		client:               client,
		namedRootHexMD5:      namedRootHexMD5,
		rootAnchorsHexSHA256: rootAnchorsHexSHA256}
}
