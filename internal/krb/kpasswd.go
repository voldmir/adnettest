package krb

import (
	"fmt"
	"net"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	KPASSWD_PORT = 464
)

type KpaswdPacket struct {
	Realm string
	SPN   string
}

func (i *KpaswdPacket) Protocols() []string {
	return []string{common.PROTO_TCP, common.PROTO_UDP}
}

func (i *KpaswdPacket) Port() string {
	return fmt.Sprint(KPASSWD_PORT)
}

func (i *KpaswdPacket) Buff() []byte {
	return make([]byte, 600)
}

func (i *KpaswdPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {

	client := new_krb_client(i.Realm, i.SPN)

	bt, err := client.GetMsgAPReqKpasswd()
	if err != nil {
		fmt.Printf(">>>>>> %s\n", err)
	}

	return bt
}
