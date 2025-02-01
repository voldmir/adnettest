package ldap

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	LDAPS_SERVICE_PORT = 636
)

type LDAPSPacket struct {
}

func (i *LDAPSPacket) Protocols() []string {
	return []string{common.PROTO_TCP}
}

func (i *LDAPSPacket) Port() string {
	return fmt.Sprint(LDAPS_SERVICE_PORT)
}

func (i *LDAPSPacket) Buff() []byte {
	return []byte{}
}

func (i *LDAPSPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {
	return []byte{}
}
