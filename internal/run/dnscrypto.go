package run

import (
	"path/filepath"

	"github.com/qdm12/updated/internal/constants"
)

func (r *runner) buildNamedRoot() error {
	// Build named root from internic.net
	namedRoot, err := r.dnscrypto.GetNamedRoot()
	if err != nil {
		return err
	}
	return r.fileManager.WriteToFile(
		filepath.Join(r.settings.OutputDir, constants.NamedRootFilename),
		namedRoot)
}

func (r *runner) buildRootAnchorsAndKeys() error {
	// Build root anchors XML from data.iana.org
	rootAnchorsXML, err := r.dnscrypto.GetRootAnchorsXML()
	if err != nil {
		return err
	}
	rootKeys, err := r.dnscrypto.ConvertRootAnchorsToRootKeys(rootAnchorsXML)
	if err != nil {
		return err
	}
	if err := r.fileManager.WriteToFile(
		filepath.Join(r.settings.OutputDir, constants.RootAnchorsFilename),
		rootAnchorsXML); err != nil {
		return err
	}
	return r.fileManager.WriteLinesToFile(
		filepath.Join(r.settings.OutputDir, constants.RootKeyFilename),
		rootKeys)
}
