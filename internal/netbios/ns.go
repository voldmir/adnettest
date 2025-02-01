package netbios

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	NS_SRVC_PORT      = 137
	NBT_QCLASS_IP     = 0x0001
	NBT_QTYPE_NETBIOS = 0x0020
)

type NSHeader struct {
	TransportID uint16
	Flags       uint16
	Qdcount     uint16
	Ancount     uint16
	Nscount     uint16
	Arcount     uint16
}

type NSQuestion struct {
	Header        NSHeader
	QuestionName  NmbName
	QuestionType  uint16
	QuestionClass uint16
}

type NBNSPacket struct {
}

func (i *NBNSPacket) Protocols() []string {
	return []string{common.PROTO_UDP}
}

func (i *NBNSPacket) Port() string {
	return fmt.Sprint(NS_SRVC_PORT)
}

func (i *NBNSPacket) Buff() []byte {
	return make([]byte, 600)
}

func (i *NBNSPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {

	var size uint16

	question := NSQuestion{
		Header: NSHeader{
			TransportID: generate_name_trn_id(),
			Flags:       0x0000,
			Qdcount:     0x0001,
		},
		QuestionName:  getNmbName("*", "", NBT_NAME_CLIENT, &size),
		QuestionType:  NBT_QTYPE_NETBIOS,
		QuestionClass: NBT_QCLASS_IP, //IN
	}

	bt, err := common.Marshal(question)
	if err != nil {
		fmt.Printf(">>>>>> %v\n", err)
	}

	return bt
}
