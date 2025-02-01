package rsync

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	RSYNC_SERVICE_PORT = 873
)

type RSYNCPacket struct {
}

func (i *RSYNCPacket) Protocols() []string {
	return []string{common.PROTO_TCP}
}

func (i *RSYNCPacket) Port() string {
	return fmt.Sprint(RSYNC_SERVICE_PORT)
}

func (i *RSYNCPacket) Buff() []byte {
	return []byte{}
}

func (i *RSYNCPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {
	return []byte{}
}
