package netbios

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	SSN_SRVC_PORT = 139
)

type NBSSNPacket struct {
}

func (i *NBSSNPacket) Protocols() []string {
	return []string{common.PROTO_TCP}
}

func (i *NBSSNPacket) Port() string {
	return fmt.Sprint(SSN_SRVC_PORT)
}

func (i *NBSSNPacket) Buff() []byte {
	return []byte{}
}

func (i *NBSSNPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {
	return []byte{}
}
