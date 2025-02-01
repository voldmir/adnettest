package dns

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	NAME_SERVICE_PORT = 53
)

type DNSPacket struct {
}

func (i *DNSPacket) Protocols() []string {
	return []string{common.PROTO_TCP, common.PROTO_UDP}
}

func (i *DNSPacket) Port() string {
	return fmt.Sprint(NAME_SERVICE_PORT)
}

func (i *DNSPacket) Buff() []byte {
	return make([]byte, 600)
}

func (i *DNSPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {

	return []byte{
		0x00, 0x01, // Transaction ID
		0x00, 0x00, // Flags: Standard query
		0x00, 0x01, // Questions
		0x00, 0x00, // Answer RRs
		0x00, 0x00, // Authority RRs
		0x00, 0x00, // Additional RRs
	}

}
