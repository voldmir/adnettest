package ldap

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	LDAPGC_SERVICE_PORT = 3268
)

type LDAPGCPacket struct {
}

func (i *LDAPGCPacket) Protocols() []string {
	return []string{common.PROTO_TCP}
}

func (i *LDAPGCPacket) Port() string {
	return fmt.Sprint(LDAPGC_SERVICE_PORT)
}

func (i *LDAPGCPacket) Buff() []byte {
	return []byte{}
}

func (i *LDAPGCPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {
	return []byte{}
}
