package smb

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	SMB_SERVICE_PORT = 445
)

type SMBPacket struct {
}

func (i *SMBPacket) Protocols() []string {
	return []string{common.PROTO_TCP}
}

func (i *SMBPacket) Port() string {
	return fmt.Sprint(SMB_SERVICE_PORT)
}

func (i *SMBPacket) Buff() []byte {
	return []byte{}
}

func (i *SMBPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {
	return []byte{}
}
