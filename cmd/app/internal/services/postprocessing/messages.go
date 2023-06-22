package postprocessing

import (
	"encoding/binary"
	"fmt"
	"go-flow2sqlite/cmd/app/internal/models"
	decodeflow "go-flow2sqlite/cmd/app/internal/services/udp"

	"github.com/sirupsen/logrus"
)

func PPToString(record *decodeflow.DecodedRecord, SubNets, IgnorList []string, ds *models.Devices) (string, models.BMes) {
	binRecord := record.BinaryRecord
	header := record.Header
	remoteAddr := record.Host
	srcmacB := make([]byte, 8)
	dstmacB := make([]byte, 8)
	binary.BigEndian.PutUint16(srcmacB, binRecord.SrcAs)
	binary.BigEndian.PutUint16(dstmacB, binRecord.DstAs)
	// srcmac = srcmac[2:8]
	// dstmac = dstmac[2:8]

	var protocol, message string
	binaryMassege := models.BMes{}

	switch fmt.Sprintf("%v", binRecord.Protocol) {
	case "6":
		protocol = "TCP_PACKET"
	case "17":
		protocol = "UDP_PACKET"
	case "1":
		protocol = "ICMP_PACKET"

	default:
		protocol = "OTHER_PACKET"
	}

	ok := CheckEntryInSubNet(intToIPv4Addr(binRecord.Ipv4DstAddrInt), SubNets)
	ok2 := CheckEntryInSubNet(intToIPv4Addr(binRecord.Ipv4SrcAddrInt), SubNets)

	if ok && !ok2 {
		ipDst := intToIPv4Addr(binRecord.Ipv4DstAddrInt).String()
		if inIgnor(ipDst, IgnorList) {
			return "", models.BMes{}
		}
		response := GetInfo(ds, ipDst, fmt.Sprint(header.UnixSec))
		message = fmt.Sprintf("%v.000 %6v %v %v/- %v HEAD %v:%v %v FIRSTUP_PARENT/%v packet_netflow/:%v %v %v",
			header.UnixSec,                       // time
			binRecord.LastInt-binRecord.FirstInt, //delay
			ipDst,                                // dst ip
			protocol,                             // protocol
			binRecord.InBytes,                    // size
			intToIPv4Addr(binRecord.Ipv4SrcAddrInt).String(), //src ip
			binRecord.L4SrcPort, // src port
			response.Mac,        // dstmac
			remoteAddr,          // routerIP
			// net.HardwareAddr(srcmacB).String(), // srcmac
			binRecord.L4DstPort, // dstport
			response.HostName,
			response.Comments,
		)
		binaryMassege = models.BMes{
			Time:       header.UnixSec,
			Delay:      binRecord.LastInt - binRecord.FirstInt,
			IPDst:      ipDst,
			Protocol:   protocol,
			SizeOfByte: binRecord.InBytes,
			IPSrc:      intToIPv4Addr(binRecord.Ipv4SrcAddrInt).String(),
			PortSRC:    binRecord.L4SrcPort,
			MacDst:     response.Mac,
			RouterIP:   remoteAddr,
			PortDst:    binRecord.L4DstPort,
		}

	} else if !ok && ok2 {
		ipDst := intToIPv4Addr(binRecord.Ipv4SrcAddrInt).String()
		if inIgnor(ipDst, IgnorList) {
			return "", models.BMes{}
		}
		response := GetInfo(ds, ipDst, fmt.Sprint(header.UnixSec))
		message = fmt.Sprintf("%v.000 %6v %v %v/- %v HEAD %v:%v %v FIRSTUP_PARENT/%v packet_netflow_inverse/:%v %v %v",
			header.UnixSec,                       // time
			binRecord.LastInt-binRecord.FirstInt, //delay
			ipDst,                                //src ip - Local
			protocol,                             // protocol
			binRecord.InBytes,                    // size
			intToIPv4Addr(binRecord.Ipv4DstAddrInt).String(), // dst ip - Inet
			binRecord.L4DstPort, // dstport
			response.Mac,        // dstmac
			remoteAddr,          // routerIP
			// net.HardwareAddr(srcmacB).String(), // srcmac
			binRecord.L4SrcPort, // src port
			response.HostName,
			response.Comments,
		)
		binaryMassege = models.BMes{
			Time:       header.UnixSec,
			Delay:      binRecord.LastInt - binRecord.FirstInt,
			IPDst:      ipDst,
			Protocol:   protocol,
			SizeOfByte: binRecord.InBytes,
			IPSrc:      intToIPv4Addr(binRecord.Ipv4SrcAddrInt).String(),
			PortSRC:    binRecord.L4SrcPort,
			MacDst:     response.Mac,
			RouterIP:   remoteAddr,
			PortDst:    binRecord.L4DstPort,
		}

	}
	return message, binaryMassege
}

func PreparingForStore(inChan chan models.BMes, outChan chan decodeflow.DecodedRecord, subNets, ignorList []string, ds *models.Devices) {
	for record := range outChan {
		logrus.Tracef("Get from outputChannel:%v", record)
		message, bm := PPToString(&record, subNets, ignorList, ds)
		logrus.Tracef("Decoded record (%v) to message (%v)", record, message)
		message = filtredMessage(message, ignorList)
		if message == "" {
			continue
		}
		inChan <- bm
	}
}

func GetInfo(ds *models.Devices, ip, time string) *models.DeviceType {
	ds.MX.RLock()
	device, ok := ds.M[ip]
	ds.MX.RUnlock()
	if !ok {
		device.Mac = ip
		device.HostName = "host_" + ip
	}
	return &device
}
