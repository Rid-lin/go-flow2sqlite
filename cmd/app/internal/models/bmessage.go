package models

type BMes struct {
	Time       uint32 // time
	Delay      uint32 //delay
	IPSrc      string //src ip - Local
	Protocol   string // protocol
	SizeOfByte uint32 // size
	IPDst      string // dst ip - Inet
	PortDst    uint16 // dstport
	MacDst     string // dstmac
	RouterIP   string // routerIP
	PortSRC    uint16 // src port

}
