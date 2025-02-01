package ldap

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	LDAPGCS_SERVICE_PORT = 3269
)

type LDAPGCSPacket struct {
}

func (i *LDAPGCSPacket) Protocols() []string {
	return []string{common.PROTO_TCP}
}

func (i *LDAPGCSPacket) Port() string {
	return fmt.Sprint(LDAPGCS_SERVICE_PORT)
}

func (i *LDAPGCSPacket) Buff() []byte {
	return []byte{}
}

func (i *LDAPGCSPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {
	return []byte{}
}
