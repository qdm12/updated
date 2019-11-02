package ips

import (
	"net"
	"sort"

	"github.com/qdm12/golibs/logging"
	"github.com/yl2chen/cidranger"
)

// CleanIPs removes duplicates IPs and CIDRs, and IPs contained in CIDRs
func CleanIPs(IPs []string) (cleanIPs []string, removedCount int) {
	uniqueIPs := makeUniqueIPs(IPs)
	uniqueCIDRs := makeUniqueCIDRs(IPs)
	ranger := buildCIDRRanger(uniqueCIDRs)
	uniqueIPs = removeIPsInCIDRs(uniqueIPs, ranger)
	// TODO CIDR inside CIDR?
	// TODO Combine CIDRs (check mask only)
	cleanIPs = sortStringSlice(append(uniqueCIDRs, uniqueIPs...))
	return cleanIPs, len(IPs) - len(cleanIPs)
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

func buildCIDRRanger(CIDRs []string) (ranger cidranger.Ranger) {
	ranger = cidranger.NewPCTrieRanger()
	for _, CIDR := range CIDRs {
		_, CIDRPtr, err := net.ParseCIDR(CIDR)
		if err == nil {
			ranger.Insert(cidranger.NewBasicRangerEntry(*CIDRPtr))
		}
	}
	return ranger
}

func removeIPsInCIDRs(IPs []string, ranger cidranger.Ranger) (cleanedIPs []string) {
	for _, IP := range IPs {
		netIP := net.ParseIP(IP)
		if netIP == nil {
			continue
		}
		contains, err := ranger.Contains(netIP)
		if err != nil {
			logging.Werr(err)
			continue
		}
		if !contains {
			cleanedIPs = append(cleanedIPs, IP)
		}
	}
	return cleanedIPs
}

func sortStringSlice(slice []string) []string {
	var sorted sort.StringSlice
	for _, element := range slice {
		sorted = append(sorted, element)
	}
	sorted.Sort()
	return sorted
}
