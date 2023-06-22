package decodeflow

import (
	"net"
)

func intToIPv4Addr(intAddr uint32) net.IP {

	return net.IPv4(
		byte(intAddr>>24),
		byte(intAddr>>16),
		byte(intAddr>>8),
		byte(intAddr))
}
