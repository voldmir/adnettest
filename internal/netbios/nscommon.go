package netbios

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"strings"
	"time"

	"github.com/voldmir/adnettest/internal/common"
)

const (
	NBT_NAME_CLIENT = 0x00
)

type DATA_BLOB struct {
	data   string
	length uint8
}

// NetBIOS DATAGRAM HEADER
type Header struct {
	MsgType    uint8
	Flags      uint8
	DgmId      uint16
	SourceIP   uint32
	SourcePort uint16
}

type NBDataGramDirectUnique struct {
	Header       Header
	DgmLength    uint16
	PacketOffset uint16
	SourceName   NmbName
	DestName     NmbName
}

type NmbName struct {
	name      string
	scope     string
	name_type string
}

func formatStringASCII(str string) string {
	return str + string(0x00)
}

func formatStringUTF8(str string) string {
	out := bytes.Buffer{}
	for _, char := range []byte(formatStringASCII(str)) {
		out.Write([]byte{char, 0x00})
	}
	return out.String()
}

func getNmbName(name string, scope string, name_type byte, size *uint16) NmbName {
	*size += uint16(2 + len(scope))

	return NmbName{
		name:      string(0x20) + netbios_name_encoding(name, size),
		scope:     scope,
		name_type: formatStringASCII(str_compressed(string(name_type), size)),
	}
}

func GetHostAddress(s net.IP) uint32 {
	b := []byte(s)
	ip := binary.BigEndian.Uint32(b)
	return ip
}

func generate_name_trn_id() uint16 {
	now := time.Now().Second()
	name_trn_id := uint16(now) % uint16(0x7FFF)
	name_trn_id += uint16(os.Getppid()) % 100
	name_trn_id = (name_trn_id + 1) % uint16(0x7FFF)
	return name_trn_id
}

func str_compressed(str string, size *uint16) string {
	l := len(str)
	var buf bytes.Buffer
	*size += uint16(l * 2)
	for i := 0; i < l; i++ {
		char := str[i]
		buf.WriteByte(((char & 0xF0) >> 4) + 0x41)
		buf.WriteByte((char & 0x0F) + 0x41)
	}
	return buf.String()

}

func netbios_name_encoding(name string, size *uint16) string {
	s := strings.ToUpper(name)
	s = common.IndentByRuneRight(s, 15, ' ')
	return str_compressed(s, size)
}
