package dnscrypto

import "github.com/qdm12/golibs/network"

type DNSCrypto interface {
	GetNamedRoot() (namedRoot []byte, err error)
	GetRootAnchorsXML() (rootAnchorsXML []byte, err error)
	ConvertRootAnchorsToRootKeys(rootAnchorsXML []byte) (rootKeys []string, err error)
}

type dnsCrypto struct {
	client network.Client
}

func NewDNSCrypto(client network.Client) DNSCrypto {
	return &dnsCrypto{client}
}
