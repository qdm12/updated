package run

import (
	"path/filepath"

	"github.com/qdm12/updated/internal/constants"
)

func (r *runner) buildBlockLists(buildHostnames, buildIps func() ([]string, error),
	hostnamesFilename, ipsFilename string) error {
	hostnames, err := buildHostnames()
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
		newIPs, err := buildIps()
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

func (r *runner) buildMalicious() error {
	return r.buildBlockLists(
		r.hostnamesBuilder.BuildMalicious,
		r.ipsBuilder.BuildMalicious,
		constants.MaliciousHostnamesFilename,
		constants.MaliciousIPsFilename,
	)
}

func (r *runner) buildAds() error {
	return r.buildBlockLists(
		r.hostnamesBuilder.BuildAds,
		nil,
		constants.MaliciousHostnamesFilename,
		constants.MaliciousIPsFilename,
	)
}

func (r *runner) buildSurveillance() error {
	return r.buildBlockLists(
		r.hostnamesBuilder.BuildSurveillance,
		nil,
		constants.SurveillanceHostnamesFilename,
		constants.SurveillanceIPsFilename,
	)
}
