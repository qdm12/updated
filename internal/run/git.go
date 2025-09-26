package run

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/updated/internal/settings"
	"github.com/qdm12/updated/pkg/git"
)

func setupGit(ctx context.Context, settings settings.Settings,
	logger logging.Logger,
) (gitUploader, error) {
	gitSettings := settings.Git

	if !*gitSettings.Enabled {
		return gitUploader{logger: logger}, nil
	}

	keyPassphrase, err := readSSHKeyPassphrase(*gitSettings.SSHKeyPassphrase)
	if err != nil {
		return gitUploader{}, fmt.Errorf("reading SSH key passphrase: %w", err)
	}

	// Setup Git repository
	client, err := git.New(
		gitSettings.SSHKnownHosts,
		gitSettings.SSHKey,
		keyPassphrase,
		gitSettings.GitURL,
		settings.OutputDir)
	if err != nil {
		return gitUploader{}, fmt.Errorf("setting up Git client: %w", err)
	}
	err = client.Pull(ctx)
	if err != nil {
		return gitUploader{}, fmt.Errorf("pulling latest changes: %w", err)
	}

	return gitUploader{
		client: client,
		logger: logger,
	}, nil
}

func readSSHKeyPassphrase(passphraseFile string) (passphrase string, err error) {
	if passphraseFile == "" {
		return "", nil
	}

	content, err := os.ReadFile(passphraseFile) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("reading SSH key passphrase file: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}

type gitUploader struct {
	client *git.Client
	logger logging.Logger
}

func (g *gitUploader) UploadAllChanges(ctx context.Context, message string) error {
	if g.client == nil {
		g.logger.Info("Git upload skipped: disabled")
		return nil
	}

	err := g.client.UploadAllChanges(ctx, message)
	if err != nil {
		return fmt.Errorf("uploading changes: %w", err)
	}
	g.logger.Info("Committed to Git: " + message)
	return nil
}
