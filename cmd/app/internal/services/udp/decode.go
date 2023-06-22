package decodeflow

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/sirupsen/logrus"
)

// NetFlow v5 implementation

type Header struct {
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

type BinaryRecord struct {
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

type DecodedRecord struct {
	Header
	BinaryRecord

	Host              string
	SamplingAlgorithm uint8
	Ipv4SrcAddr       string
	Ipv4DstAddr       string
	Ipv4NextHop       string
	SrcHostName       string
	DstHostName       string
	Duration          uint16
}

func HandlePacket(buf *bytes.Buffer, remoteAddr *net.UDPAddr, outputChannel chan DecodedRecord) {
	header := Header{}
	err := binary.Read(buf, binary.BigEndian, &header)
	if err != nil {
		logrus.Printf("Error: %v\n", err)
	} else {

		for i := 0; i < int(header.FlowRecords); i++ {
			record := BinaryRecord{}
			err := binary.Read(buf, binary.BigEndian, &record)
			if err != nil {
				logrus.Printf("binary.Read failed: %v\n", err)
				break
			}

			decodedRecord := decodeRecord(&header, &record, remoteAddr)
			logrus.Tracef("Send to outputChannel:%v", decodedRecord)
			outputChannel <- decodedRecord
		}
	}
}

func decodeRecord(header *Header, binRecord *BinaryRecord, remoteAddr *net.UDPAddr) DecodedRecord {

	decodedRecord := DecodedRecord{

		Host: remoteAddr.IP.String(),

		Header: *header,

		BinaryRecord: *binRecord,

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
