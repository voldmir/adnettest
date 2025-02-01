package krb

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	KERBEROS_PORT = 88
)

type KerberosPacket struct {
	Realm string
	SPN   string
}

func (i *KerberosPacket) Protocols() []string {
	return []string{common.PROTO_TCP, common.PROTO_UDP}
}

func (i *KerberosPacket) Port() string {
	return fmt.Sprint(KERBEROS_PORT)
}

func (i *KerberosPacket) Buff() []byte {
	return make([]byte, 600)
}

func (i *KerberosPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {

	client := new_krb_client(i.Realm, i.SPN)

	bt, err := client.GetMsgASReq()
	if err != nil {
		fmt.Printf(">>>>>> %s\n", err)
	}

	return bt
}
