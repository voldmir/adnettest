package www

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	HTTPS_SERVICE_PORT = 443
)

type HTTPSPacket struct {
}

func (i *HTTPSPacket) Protocols() []string {
	return []string{common.PROTO_TCP}
}

func (i *HTTPSPacket) Port() string {
	return fmt.Sprint(HTTPS_SERVICE_PORT)
}

func (i *HTTPSPacket) Buff() []byte {
	return []byte{}
}

func (i *HTTPSPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {
	return []byte{}
}
