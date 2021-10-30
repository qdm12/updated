package run

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/qdm12/updated/internal/constants"
)

func (r *runner) buildBlockLists(ctx context.Context, buildHostnames,
	buildIps func(ctx context.Context) ([]string, error),
	hostnamesFilename, ipsFilename string) error {
	hostnames, err := buildHostnames(ctx)
	if err != nil {
		return err
	}

	hostnamesFilepath := filepath.Join(r.settings.OutputDir, hostnamesFilename)
	file, err := os.OpenFile(hostnamesFilepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.WriteString(strings.Join(hostnames, "\n"))
	if err != nil {
		_ = file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
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
	r.logger.Info("Trimmed down %d IP address lines", removedCount)

	ipsFilepath := filepath.Join(r.settings.OutputDir, ipsFilename)
	file, err = os.OpenFile(ipsFilepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.WriteString(strings.Join(IPs, "\n"))
	if err != nil {
		_ = file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func (r *runner) buildMalicious(ctx context.Context) error {
	return r.buildBlockLists(ctx,
		r.hostnamesBuilder.BuildMalicious,
		r.ipsBuilder.BuildMalicious,
		constants.MaliciousHostnamesFilename,
		constants.MaliciousIPsFilename,
	)
}

func (r *runner) buildAds(ctx context.Context) error {
	return r.buildBlockLists(ctx,
		r.hostnamesBuilder.BuildAds,
		nil,
		constants.AdsHostnamesFilename,
		constants.AdsIPsFilename,
	)
}

func (r *runner) buildSurveillance(ctx context.Context) error {
	return r.buildBlockLists(ctx,
		r.hostnamesBuilder.BuildSurveillance,
		nil,
		constants.SurveillanceHostnamesFilename,
		constants.SurveillanceIPsFilename,
	)
}
