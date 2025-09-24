package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/qdm12/updated/internal/constants"
)

func (r *Runner) buildBlockLists(ctx context.Context, buildHostnames,
	buildIps func(ctx context.Context) ([]string, error),
	hostnamesFilename, ipsFilename string,
) error {
	hostnames, err := buildHostnames(ctx)
	if err != nil {
		return err
	}

	hostnamesFilepath := filepath.Join(r.settings.OutputDir, hostnamesFilename)
	err = writeLines(hostnamesFilepath, hostnames)
	if err != nil {
		return fmt.Errorf("writing hostnames: %w", err)
	}

	IPs := []string{}
	if r.settings.ResolveHostnames {
		IPs = append(IPs, r.ipsBuilder.BuildIPsFromHostnames(hostnames)...)
	}
	if buildIps != nil {
		newIPs, err := buildIps(ctx)
		if err != nil {
			return err
		}
		IPs = append(IPs, newIPs...)
	}
	IPs, removedCount, warnings := r.ipsBuilder.CleanIPs(IPs)
	for _, warning := range warnings {
		r.logger.Warn(warning)
	}
	r.logger.Info(fmt.Sprintf("Trimmed down %d IP address lines", removedCount))

	ipsFilepath := filepath.Join(r.settings.OutputDir, ipsFilename)
	err = writeLines(ipsFilepath, IPs)
	if err != nil {
		return fmt.Errorf("writing IPs: %w", err)
	}

	return nil
}

func writeLines(filePath string, lines []string) error {
	const perms = 0o600
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perms) //nolint:gosec
	if err != nil {
		return err
	}

	_, err = file.WriteString(strings.Join(lines, "\n"))
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

func (r *Runner) buildMalicious(ctx context.Context) error {
	return r.buildBlockLists(ctx,
		r.hostnamesBuilder.BuildMalicious,
		r.ipsBuilder.BuildMalicious,
		constants.MaliciousHostnamesFilename,
		constants.MaliciousIPsFilename,
	)
}

func (r *Runner) buildAds(ctx context.Context) error {
	return r.buildBlockLists(ctx,
		r.hostnamesBuilder.BuildAds,
		nil,
		constants.AdsHostnamesFilename,
		constants.AdsIPsFilename,
	)
}

func (r *Runner) buildSurveillance(ctx context.Context) error {
	return r.buildBlockLists(ctx,
		r.hostnamesBuilder.BuildSurveillance,
		nil,
		constants.SurveillanceHostnamesFilename,
		constants.SurveillanceIPsFilename,
	)
}
