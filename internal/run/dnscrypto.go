package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qdm12/updated/internal/constants"
)

func (r *Runner) buildNamedRoot(ctx context.Context) error {
	// Build named root from internic.net
	namedRoot, err := r.dnscrypto.DownloadNamedRoot(ctx)
	if err != nil {
		return fmt.Errorf("downloading named root: %w", err)
	}

	filepath := filepath.Join(r.settings.OutputDir, constants.NamedRootFilename)
	const perms = 0o600
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perms) //nolint:gosec
	if err != nil {
		return err
	}

	_, err = file.Write(namedRoot)
	if err != nil {
		_ = file.Close()
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) buildRootAnchorsAndKeys(ctx context.Context) error {
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
	err = writeFile(xmlFilepath, rootAnchorsXML)
	if err != nil {
		return fmt.Errorf("writing root anchors XML: %w", err)
	}

	rootKeysFilepath := filepath.Join(r.settings.OutputDir, constants.RootKeyFilename)
	err = writeLines(rootKeysFilepath, rootKeys)
	if err != nil {
		return fmt.Errorf("writing root keys: %w", err)
	}

	return nil
}

func writeFile(filePath string, data []byte) error {
	const perms = 0o600

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perms) //nolint:gosec
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		_ = file.Close()
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
