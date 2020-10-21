package run

import (
	"context"
	"path/filepath"

	"github.com/qdm12/updated/internal/constants"
)

func (r *runner) buildBlockLists(ctx context.Context, buildHostnames,
	buildIps func(ctx context.Context) ([]string, error),
	hostnamesFilename, ipsFilename string) error {
	hostnames, err := buildHostnames(ctx)
	if err != nil {
		return err
	}
	if err := r.fileManager.WriteLinesToFile(
		filepath.Join(r.settings.OutputDir, hostnamesFilename),
		hostnames); err != nil {
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
	return r.fileManager.WriteLinesToFile(
		filepath.Join(r.settings.OutputDir, ipsFilename),
		IPs)
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
