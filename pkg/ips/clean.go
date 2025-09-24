package ips

import (
	"net"
	"sort"

	"github.com/yl2chen/cidranger"
)

// CleanIPs removes duplicates IPs and CIDRs, and IPs contained in CIDRs.
func (b *Builder) CleanIPs(ips []string) (cleanIPs []string, removedCount int, warnings []string) {
	uniqueIPs := makeUniqueIPs(ips)
	uniqueCIDRs := makeUniqueCIDRs(ips)
	ranger, err := buildCIDRRanger(uniqueCIDRs)
	if err != nil {
		warnings = append(warnings, err.Error())
	}
	uniqueIPs, newWarnings := removeIPsInCIDRs(uniqueIPs, ranger)
	warnings = append(warnings, newWarnings...)
	// TODO CIDR inside CIDR?
	// TODO Combine CIDRs (check mask only)
	cleanIPs = sortStringSlice(append(uniqueCIDRs, uniqueIPs...))
	return cleanIPs, len(ips) - len(cleanIPs), warnings
}

func makeUniqueIPs(lines []string) (uniqueIPs []string) {
	IPMap := make(map[string]struct{})
	for _, line := range lines {
		// only process single IP
		if IP := net.ParseIP(line); IP != nil {
			IPMap[IP.String()] = struct{}{}
		}
	}
	for IP := range IPMap {
		uniqueIPs = append(uniqueIPs, IP)
	}
	return uniqueIPs
}

func makeUniqueCIDRs(lines []string) (uniqueCIDRs []string) {
	CIDRMap := make(map[string]struct{})
	for _, line := range lines {
		// only process CIDR
		_, CIDRPtr, err := net.ParseCIDR(line)
		if err == nil {
			CIDRMap[CIDRPtr.String()] = struct{}{}
			continue
		}
	}
	for CIDR := range CIDRMap {
		uniqueCIDRs = append(uniqueCIDRs, CIDR)
	}
	return uniqueCIDRs
}

//nolint:ireturn
func buildCIDRRanger(cidrs []string) (ranger cidranger.Ranger, err error) {
	ranger = cidranger.NewPCTrieRanger()
	for _, CIDR := range cidrs {
		_, CIDRPtr, err := net.ParseCIDR(CIDR)
		if err == nil {
			err = ranger.Insert(cidranger.NewBasicRangerEntry(*CIDRPtr))
			if err != nil {
				return nil, err
			}
		}
	}
	return ranger, nil
}

func removeIPsInCIDRs(ips []string, ranger cidranger.Ranger) (cleanedIPs []string, warnings []string) {
	for _, IP := range ips {
		netIP := net.ParseIP(IP)
		if netIP == nil {
			continue
		}
		contains, err := ranger.Contains(netIP)
		if err != nil {
			warnings = append(warnings, err.Error())
			continue
		}
		if !contains {
			cleanedIPs = append(cleanedIPs, IP)
		}
	}
	return cleanedIPs, warnings
}

func sortStringSlice(slice []string) []string {
	var sorted sort.StringSlice
	for _, element := range slice {
		sorted = append(sorted, element)
	}
	sorted.Sort()
	return sorted
}
