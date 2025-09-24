package dnscrypto

const (
	// NamedRootMD5Sum is the default expected MD5 checksum for the named root file.
	NamedRootMD5Sum = "076cfeb40394314adf28b7be79e6ecb1"
	// RootAnchorsSHA256Sum is the default expected SHA256 checksum for the root anchors file.
	RootAnchorsSHA256Sum = "45336725f9126db810a59896ae93819de743c416262f79c4444042c92e520770"
)

// SetNamedRootHexMD5 sets the expected MD5 checksum for the named root file.
func (d *DNSCrypto) SetNamedRootHexMD5(namedRootHexMD5 string) {
	d.namedRootHexMD5Mu.Lock()
	defer d.namedRootHexMD5Mu.Unlock()
	d.namedRootHexMD5 = namedRootHexMD5
}

// SetRootAnchorsHexSHA256 sets the expected SHA256 checksum for the root anchors file.
func (d *DNSCrypto) SetRootAnchorsHexSHA256(rootAnchorsHexSHA256 string) {
	d.rootAnchorsHexSHA256Mu.Lock()
	defer d.rootAnchorsHexSHA256Mu.Unlock()
	d.rootAnchorsHexSHA256 = rootAnchorsHexSHA256
}

func (d *DNSCrypto) getNamedRootHexMD5() string {
	d.namedRootHexMD5Mu.RLock()
	defer d.namedRootHexMD5Mu.RUnlock()
	return d.namedRootHexMD5
}

func (d *DNSCrypto) getRootAnchorsHexSHA256() string {
	d.rootAnchorsHexSHA256Mu.RLock()
	defer d.rootAnchorsHexSHA256Mu.RUnlock()
	return d.rootAnchorsHexSHA256
}
