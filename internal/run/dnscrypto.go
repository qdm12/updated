package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/qdm12/updated/internal/constants"
)

func (r *runner) buildNamedRoot(ctx context.Context) error {
	// Build named root from internic.net
	namedRoot, err := r.dnscrypto.DownloadNamedRoot(ctx)
	if err != nil {
		return fmt.Errorf("downloading named root: %w", err)
	}

	filepath := filepath.Join(r.settings.OutputDir, constants.NamedRootFilename)
	file, err := r.osOpenFile(filepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write(namedRoot)
	if err != nil {
		_ = file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func (r *runner) buildRootAnchorsAndKeys(ctx context.Context) error {
	// Build root anchors XML from data.iana.org
	rootAnchorsXML, err := r.dnscrypto.DownloadRootAnchorsXML(ctx)
	if err != nil {
		return fmt.Errorf("downloading root anchors XML: %w", err)
	}
	rootKeys, err := r.dnscrypto.ConvertRootAnchorsToRootKeys(rootAnchorsXML)
	if err != nil {
		return err
	}

	xmlFilepath := filepath.Join(r.settings.OutputDir, constants.RootAnchorsFilename)
	file, err := r.osOpenFile(xmlFilepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write(rootAnchorsXML)
	if err != nil {
		_ = file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	rootKeysFilepath := filepath.Join(r.settings.OutputDir, constants.RootKeyFilename)
	file, err = r.osOpenFile(rootKeysFilepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.WriteString(strings.Join(rootKeys, "\n"))
	if err != nil {
		_ = file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}
