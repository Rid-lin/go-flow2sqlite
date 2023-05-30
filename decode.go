package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

// NetFlow v5 implementation

type header struct {
	Version          uint16
	FlowRecords      uint16
	Uptime           uint32
	UnixSec          uint32
	UnixNsec         uint32
	FlowSeqNum       uint32
	EngineType       uint8
	EngineID         uint8
	SamplingInterval uint16
}

type binaryRecord struct {
	Ipv4SrcAddrInt uint32
	Ipv4DstAddrInt uint32
	Ipv4NextHopInt uint32
	InputSnmp      uint16
	OutputSnmp     uint16
	InPkts         uint32
	InBytes        uint32
	FirstInt       uint32
	LastInt        uint32
	L4SrcPort      uint16
	L4DstPort      uint16
	_              uint8
	TCPFlags       uint8
	Protocol       uint8
	SrcTos         uint8
	SrcAs          uint16
	DstAs          uint16
	SrcMask        uint8
	DstMask        uint8
	_              uint16
}

type decodedRecord struct {
	header
	binaryRecord

	Host              string
	SamplingAlgorithm uint8
	Ipv4SrcAddr       string
	Ipv4DstAddr       string
	Ipv4NextHop       string
	SrcHostName       string
	DstHostName       string
	Duration          uint16
}

type bmes struct {
	time       uint32 // time
	delay      uint32 //delay
	ipSrc      string //src ip - Local
	protocol   string // protocol
	size_fByte uint32 // size
	ipDst      string // dst ip - Inet
	portDst    uint16 // dstport
	macDst     string // dstmac
	routerIP   string // routerIP
	portSRC    uint16 // src port

}

func intToIPv4Addr(intAddr uint32) net.IP {

	return net.IPv4(
		byte(intAddr>>24),
		byte(intAddr>>16),
		byte(intAddr>>8),
		byte(intAddr))
}

func handlePacket(buf *bytes.Buffer, remoteAddr *net.UDPAddr, outputChannel chan decodedRecord, cfg *Config) {
	header := header{}
	err := binary.Read(buf, binary.BigEndian, &header)
	if err != nil {
		logrus.Printf("Error: %v\n", err)
	} else {

		for i := 0; i < int(header.FlowRecords); i++ {
			record := binaryRecord{}
			err := binary.Read(buf, binary.BigEndian, &record)
			if err != nil {
				logrus.Printf("binary.Read failed: %v\n", err)
				break
			}

			decodedRecord := decodeRecord(&header, &record, remoteAddr, cfg)
			logrus.Tracef("Send to outputChannel:%v", decodedRecord)
			outputChannel <- decodedRecord
		}
	}
}

func decodeRecord(header *header, binRecord *binaryRecord, remoteAddr *net.UDPAddr, cfg *Config) decodedRecord {

	decodedRecord := decodedRecord{

		Host: remoteAddr.IP.String(),

		header: *header,

		binaryRecord: *binRecord,

		Ipv4SrcAddr: intToIPv4Addr(binRecord.Ipv4SrcAddrInt).String(),
		Ipv4DstAddr: intToIPv4Addr(binRecord.Ipv4DstAddrInt).String(),
		Ipv4NextHop: intToIPv4Addr(binRecord.Ipv4NextHopInt).String(),
		Duration:    uint16((binRecord.LastInt - binRecord.FirstInt) / 1000),
	}

	// decode sampling info
	decodedRecord.SamplingAlgorithm = uint8(0x3 & (decodedRecord.SamplingInterval >> 14))
	decodedRecord.SamplingInterval = 0x3fff & decodedRecord.SamplingInterval

	return decodedRecord
}

func (t *Transport) decodeRecordToSquid(record *decodedRecord, cfg *Config) (string, bmes) {
	binRecord := record.binaryRecord
	header := record.header
	remoteAddr := record.Host
	srcmacB := make([]byte, 8)
	dstmacB := make([]byte, 8)
	binary.BigEndian.PutUint16(srcmacB, binRecord.SrcAs)
	binary.BigEndian.PutUint16(dstmacB, binRecord.DstAs)
	// srcmac = srcmac[2:8]
	// dstmac = dstmac[2:8]

	var protocol, message string
	binaryMassege := bmes{}

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

	ok := CheckEntryInSubNet(intToIPv4Addr(binRecord.Ipv4DstAddrInt), cfg.SubNets)
	ok2 := CheckEntryInSubNet(intToIPv4Addr(binRecord.Ipv4SrcAddrInt), cfg.SubNets)

	if ok && !ok2 {
		ipDst := intToIPv4Addr(binRecord.Ipv4DstAddrInt).String()
		if inIgnor(ipDst, cfg.IgnorList) {
			return "", bmes{}
		}
		response := t.GetInfo(ipDst, fmt.Sprint(header.UnixSec))
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
		binaryMassege = bmes{
			time:       header.UnixSec,
			delay:      binRecord.LastInt - binRecord.FirstInt,
			ipDst:      ipDst,
			protocol:   protocol,
			size_fByte: binRecord.InBytes,
			ipSrc:      intToIPv4Addr(binRecord.Ipv4SrcAddrInt).String(),
			portSRC:    binRecord.L4SrcPort,
			macDst:     response.Mac,
			routerIP:   remoteAddr,
			portDst:    binRecord.L4DstPort,
		}

	} else if !ok && ok2 {
		ipDst := intToIPv4Addr(binRecord.Ipv4SrcAddrInt).String()
		if inIgnor(ipDst, cfg.IgnorList) {
			return "", bmes{}
		}
		response := t.GetInfo(ipDst, fmt.Sprint(header.UnixSec))
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
		binaryMassege = bmes{
			time:       header.UnixSec,
			delay:      binRecord.LastInt - binRecord.FirstInt,
			ipDst:      ipDst,
			protocol:   protocol,
			size_fByte: binRecord.InBytes,
			ipSrc:      intToIPv4Addr(binRecord.Ipv4SrcAddrInt).String(),
			portSRC:    binRecord.L4SrcPort,
			macDst:     response.Mac,
			routerIP:   remoteAddr,
			portDst:    binRecord.L4DstPort,
		}

	}
	return message, binaryMassege
}

func (t *Transport) pipeOutputToStdoutForSquid(outputChannel chan decodedRecord, cfg *Config) {
	for record := range outputChannel {
		logrus.Tracef("Get from outputChannel:%v", record)
		message, bm := t.decodeRecordToSquid(&record, cfg)
		logrus.Tracef("Decoded record (%v) to message (%v)", record, message)
		message = filtredMessage(message, cfg.IgnorList)
		if message == "" {
			continue
		}
		// if _, err := t.fileDestination.WriteString(message + "\n"); err != nil {
		if err := t.writeMessageToStore(bm); err != nil {
			logrus.Errorf("Error writing data buffer:%v", err)
		} else {
			logrus.Tracef("Added to log:%v", message)
		}
	}
}

func (t *Transport) writeMessageToStore(bm bmes) error {
	// TODO тут должна быть вставка SQL
	ctx := context.Background()
	t.DB.ExecContext(ctx, "INSERT INTO raw_log VALUES(?,?,?,?,?,?,?,?,?,?)",
		bm.time,
		bm.delay,
		bm.ipDst,
		bm.protocol,
		bm.size_fByte,
		bm.ipSrc,
		bm.portSRC,
		bm.macDst,
		bm.routerIP,
		bm.portDst,
	)
	return nil
}

// type cacheRecord struct {
// 	Hostname string
// 	// timeout  time.Time
// }

// type Cache struct {
// 	cache map[string]cacheRecord
// 	sync.RWMutex
// }

var (
	fileDestination *os.File
)
