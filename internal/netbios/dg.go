package netbios

import (
	"fmt"
	"net"
	"strings"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	DGM_SRVC_PORT         = 138
	DGRAM_SMB             = 0xff534d42
	SMB_TRANSACTION       = 0x25
	LOGON_PRIMARY_QUERY   = 7
	NETLOGON_NT_VERSION_1 = 1
	DGRAM_DIRECT_UNIQUE   = 0x10
	DGRAM_FLAG_FIRST      = 0x02
	DGRAM_NODE_P          = 0x04
	NBT_NAME_LOGON        = 0x1C

	LOGON_SAM_LOGON_REQUEST = 18
	NBT_MAILSLOT_NETLOGON   = "\\MAILSLOT\\NET\\NETLOGON"
)

type dom_sid struct {
	sid_rev_num uint8    //Revision (1 byte): (1 byte): An 8-bit unsigned integer that specifies the revision level of the SID. This value MUST be set to 0x01.
	num_auths   uint8    // SubAuthorityCount (1 byte): An 8-bit unsigned integer that specifies the number of elements in the SubAuthority array. The maximum number of elements allowed is 15.
	id_auth     [6]uint8 // IdentifierAuthority (6 bytes)  Authority value {0,0,0,0,0,5} denotes SIDs created by the NT SID authority.
	sub_auths   []uint32 //SubAuthority (variable): A variable length array of unsigned 32-bit
}

type NETLOGON_SAM_LOGON_REQUEST struct {
	request_count uint16
	computer_name string
	user_name     string
	mailslot_name string
	acct_control  uint32 `tag:"be"`
	sid_size      uint32 `tag:"be"`
	_pad          DATA_BLOB
	sid           []byte
	nt_version    uint32 `tag:"be"`
	lmnt_token    uint16
	lm20_token    uint16
}

type smb_trans_body struct {
	wct               uint8 //[range(17,17),value(17)]
	total_param_count uint16
	total_data_count  uint16 `tag:"be"`
	max_param_count   uint16
	max_data_count    uint16
	max_setup_count   uint8
	pad               uint8
	trans_flags       uint16
	timeout           uint32 `tag:"be"`
	reserved          uint16
	param_count       uint16
	param_offset      uint16
	data_count        uint16 `tag:"be"`
	data_offset       uint16 `tag:"be"`
	setup_count       uint8  //[range(3,3),value(3)]
	pad2              uint8
	opcode            uint16 `tag:"be"`
	priority          uint16 `tag:"be"`
	_class            uint16 `tag:"be"`
	byte_count        uint16 `tag:"be"` //[value(strlen(mailslot_name)+1+data.length)]
	mailslot_name     string
	data              nbt_netlogon_packet
}

type nbt_netlogon_packet struct {
	netlogon_command uint16 `tag:"be"`
	req              []byte
}

type smb_body struct {
	trans smb_trans_body
}

type dgram_smb_packet struct {
	smb_command uint8
	err_class   uint8
	pad         uint8
	err_code    uint16
	flags       uint8
	flags2      uint16
	pid_high    uint16
	signature   uint64
	reserved    uint16
	tid         uint16
	pid         uint16
	vuid        uint16
	mid         uint16
	body        smb_body
}

type dgram_message_body struct {
	smb dgram_smb_packet
}

type dgram_message struct {
	length          uint16
	offset          uint16
	source_name     NmbName
	dest_name       NmbName
	dgram_body_type uint32
	body            dgram_message_body
}

type nbt_dgram_packet struct {
	Header Header
	data   dgram_message
}

type NBDataGramQuery struct {
	Header   Header
	DestName NmbName `tag:"optional"`
}

type NBDgramPacket struct {
	ComputerName string
	GroupName    string
}

func (i *NBDgramPacket) Protocols() []string {
	return []string{common.PROTO_UDP}
}

func (i *NBDgramPacket) Port() string {
	return fmt.Sprint(DGM_SRVC_PORT)
}

func (i *NBDgramPacket) Buff() []byte {
	return make([]byte, 600)
}

func (i *NBDgramPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {
	var bt []byte
	bt = i.netlogon_sam_logon(remoteAddr, localAddr, localPort)
	return bt
}

func (i *NBDgramPacket) netlogon_sam_logon(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {
	var size uint16
	var err error

	var msg *dgram_message
	var smb *dgram_smb_packet
	var trans *smb_trans_body
	logon := NETLOGON_SAM_LOGON_REQUEST{}

	packet := nbt_dgram_packet{}

	msg = &packet.data
	msg.length = uint16(138 + len(NBT_MAILSLOT_NETLOGON))
	msg.offset = 0

	msg.source_name = getNmbName(i.ComputerName, "", NBT_NAME_CLIENT, &size)
	msg.dest_name = getNmbName(i.GroupName, "", NBT_NAME_LOGON, &size)
	msg.dgram_body_type = DGRAM_SMB

	smb = &msg.body.smb
	smb.smb_command = SMB_TRANSACTION

	trans = &smb.body.trans
	trans.mailslot_name = formatStringASCII(NBT_MAILSLOT_NETLOGON)
	trans.data.netlogon_command = LOGON_SAM_LOGON_REQUEST

	logon.request_count = 0
	logon.computer_name = formatStringUTF8(strings.ToUpper(i.ComputerName))
	logon.user_name = formatStringUTF8(strings.ToUpper(i.ComputerName) + "$")
	logon.sid_size = 0
	logon.mailslot_name = trans.mailslot_name
	logon.nt_version = NETLOGON_NT_VERSION_1
	logon.lmnt_token = 0xFFFF
	logon.lm20_token = 0xFFFF

	trans.data.req, err = common.Marshal(logon)
	if err != nil {
		fmt.Printf(">>>>>> %v\n", err)
	}

	logon_length := uint16(len(trans.data.req) + 2)

	trans.wct = 17
	trans.total_data_count = logon_length
	trans.timeout = 1000
	trans.data_count = trans.total_data_count
	trans.data_offset = uint16(70 + len(NBT_MAILSLOT_NETLOGON))
	trans.opcode = 1
	trans.priority = 1
	trans._class = 2
	trans.setup_count = 3
	trans.byte_count = uint16(len(NBT_MAILSLOT_NETLOGON)) + logon_length + 2

	packet.Header.MsgType = DGRAM_DIRECT_UNIQUE
	packet.Header.Flags = DGRAM_FLAG_FIRST | DGRAM_NODE_P
	packet.Header.DgmId = generate_name_trn_id()
	packet.Header.SourceIP = GetHostAddress(localAddr)
	packet.Header.SourcePort = uint16(localPort)

	bt, err := common.Marshal(packet)
	if err != nil {
		fmt.Printf(">>>>>> %v\n", err)
	}

	return bt
}
