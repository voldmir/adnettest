package krb

import (
	"encoding/hex"

	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/kadmin"
	"github.com/jcmturner/gokrb5/v8/messages"
	"github.com/jcmturner/gokrb5/v8/types"
)

const (
	KRB_NT_PRINCIPAL = 1
	KRB_NT_SRV_INST  = 2
	SessionKeyHex    = "3029a003020112a12204205b7eff528c3f37b720fe681baefce8e020fe681baefce8e020fe681baefce8e0"
	TicketHex        = "613e303ca003020105a1091b0741524d2e4c4f43a21d301ba003020101a11430121b066b61646d696e1b086368616e67657077a30b3009a003020100a2020400"
	TicketEncPartHex = "302ea003020112a103020101a22204205b7e7eff7e527e8c7e3f7e377eb77e207e3f7e377eb77e207e3f7e377eb77e20"
)

type krb_client struct {
	Realm  string
	CName  types.PrincipalName
	Config config.Config
}

func (a *krb_client) GetMsgAPReqKpasswd() ([]byte, error) {
	var b []byte
	var err error
	var tkt messages.Ticket
	var eKey types.EncryptionKey
	var r kadmin.Request

	b, _ = hex.DecodeString(TicketHex)
	err = tkt.Unmarshal(b)
	if err != nil {
		goto END
	}

	b, _ = hex.DecodeString(TicketEncPartHex)
	err = tkt.EncPart.Unmarshal(b)
	if err != nil {
		goto END
	}

	b, _ = hex.DecodeString(SessionKeyHex)
	eKey.Unmarshal(b)
	if err != nil {
		goto END
	}

	r, _, err = kadmin.ChangePasswdMsg(a.CName, a.Realm, "newPasswd", tkt, eKey)
	if err != nil {
		goto END
	}

	b, err = r.Marshal()
	if err != nil {
		goto END
	}

END:
	return b, err

}

func (a *krb_client) GetMsgASReq() ([]byte, error) {
	var b []byte

	asreq, err := messages.NewASReqForTGT(a.Realm, &a.Config, a.CName)
	if err != nil {
		goto END
	}

	b, err = asreq.Marshal()
	if err != nil {
		goto END
	}

END:
	return b, err
}

func new_krb_client(realm string, spn string) krb_client {
	return krb_client{
		Realm:  realm,
		CName:  types.NewPrincipalName(KRB_NT_PRINCIPAL, spn),
		Config: config.Config{},
	}
}
