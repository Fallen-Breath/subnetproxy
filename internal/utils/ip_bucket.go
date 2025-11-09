// Package utils, ip_bucket.go
// From https://github.com/Fallen-Breath/pavonis/blob/c763d2be3f7b7f3dddcea575b01a375c7053e4de/internal/utils/ip_bucket.go
// Licensed under GPL-3.0
package utils

import "net"

func GetBucketForIP(ip net.IP) string {
	if ip == nil {
		return "unknown$"
	}
	if ip.To4() != nil {
		return "ipv4$" + ip.String()
	}
	ipBytes := ip.To16()
	if ipBytes == nil {
		return "unknown$"
	}
	// by its /64 address
	subnetIP := ip.Mask(net.CIDRMask(64, 128))
	if subnetIP == nil {
		return "unknown$"
	}
	return "ipv6$" + subnetIP.String()
}
