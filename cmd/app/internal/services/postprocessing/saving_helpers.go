package postprocessing

import (
	"net"
	"strings"

	"github.com/sirupsen/logrus"
)

func CheckEntryInSubNet(ipv4addr net.IP, subNets []string) bool {
	for _, subNet := range subNets {
		ok, err := checkIP(subNet, ipv4addr)
		if err != nil { // если ошибка, то следующая строка
			logrus.Error("Error while determining the IP subnet address:", err)
			return false

		}
		if ok {
			return true
		}
	}

	return false
}

func checkIP(subnet string, ipv4addr net.IP) (bool, error) {
	_, netA, err := net.ParseCIDR(subnet)
	if err != nil {
		return false, err
	}

	return netA.Contains(ipv4addr), nil
}

func filtredMessage(message string, IgnorList []string) string {
	for _, ignorStr := range IgnorList {
		if strings.Contains(message, ignorStr) {
			logrus.Tracef("Line of log :%v contains ignorstr:%v, skipping...", message, ignorStr)
			return ""
		}
	}
	return message
}

func inIgnor(message string, IgnorList []string) bool {
	for _, ignorStr := range IgnorList {
		if strings.Contains(message, ignorStr) {
			logrus.Tracef("Line of log :%v contains ignorstr:%v, skipping...", message, ignorStr)
			return true
		}
	}
	return false
}

func intToIPv4Addr(intAddr uint32) net.IP {

	return net.IPv4(
		byte(intAddr>>24),
		byte(intAddr>>16),
		byte(intAddr>>8),
		byte(intAddr))
}
