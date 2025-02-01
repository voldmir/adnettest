package www

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	HTTP_SERVICE_PORT = 80
)

type HTTPPacket struct {
}

func (i *HTTPPacket) Protocols() []string {
	return []string{common.PROTO_TCP}
}

func (i *HTTPPacket) Port() string {
	return fmt.Sprint(HTTP_SERVICE_PORT)
}

func (i *HTTPPacket) Buff() []byte {
	return []byte{}
}

func (i *HTTPPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {
	return []byte{}
}
