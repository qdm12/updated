package dnscrypto

const (
	NamedRootMD5Sum      = "9cdf6a12db2dc558e299601015c37948"
	RootAnchorsSHA256Sum = "45336725f9126db810a59896ae93819de743c416262f79c4444042c92e520770"
)

func (d *dnsCrypto) SetNamedRootHexMD5(namedRootHexMD5 string) {
	d.namedRootHexMD5Mu.Lock()
	defer d.namedRootHexMD5Mu.Unlock()
	d.namedRootHexMD5 = namedRootHexMD5
}

func (d *dnsCrypto) SetRootAnchorsHexSHA256(rootAnchorsHexSHA256 string) {
	d.rootAnchorsHexSHA256Mu.Lock()
	defer d.rootAnchorsHexSHA256Mu.Unlock()
	d.rootAnchorsHexSHA256 = rootAnchorsHexSHA256
}

func (d *dnsCrypto) getNamedRootHexMD5() string {
	d.namedRootHexMD5Mu.RLock()
	defer d.namedRootHexMD5Mu.RUnlock()
	return d.namedRootHexMD5
}

func (d *dnsCrypto) getRootAnchorsHexSHA256() string {
	d.rootAnchorsHexSHA256Mu.RLock()
	defer d.rootAnchorsHexSHA256Mu.RUnlock()
	return d.rootAnchorsHexSHA256
}
