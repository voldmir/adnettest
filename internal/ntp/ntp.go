package ntp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/voldmir/adnettest/internal/common"
)

type ntpTimeShort uint32
type ntpTime uint64

const NAME_SERVICE_PORT = 123

type HeaderNTP struct {
	LiVnMode       uint8 // Leap Indicator (2) + Version (3) + Mode (3)
	Stratum        uint8
	Poll           uint8
	Precision      uint8
	RootDelay      ntpTimeShort
	RootDispersion ntpTimeShort
	ReferenceID    uint32 // KoD code if Stratum == 0
	ReferenceTime  ntpTime
	OriginTime     ntpTime
	ReceiveTime    ntpTime
	TransmitTime   ntpTime
}

type NTPPacket struct {
}

func (i *NTPPacket) Protocols() []string {
	return []string{common.PROTO_UDP}
}

func (i *NTPPacket) Port() string {
	return fmt.Sprint(NAME_SERVICE_PORT)
}

func (i *NTPPacket) Buff() []byte {
	return make([]byte, 600)
}

func toNtpTime(t time.Time) ntpTime {
	const nanoPerSec = 1000000000
	nsec := uint64(t.Sub(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)))
	sec := nsec / nanoPerSec
	nsec = uint64(nsec-sec*nanoPerSec) << 32
	frac := uint64(nsec / nanoPerSec)
	if nsec%nanoPerSec >= nanoPerSec/2 {
		frac++
	}
	return ntpTime(sec<<32 | frac)
}

func (i *NTPPacket) Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte {
	var b []byte

	xmitHdr := HeaderNTP{
		LiVnMode:       0xe3,
		Poll:           0x3,
		Precision:      0xfa,
		RootDelay:      0x10000,
		RootDispersion: 0x10000,
		TransmitTime:   toNtpTime(time.Now().UTC()),
	}

	var xmitBuf bytes.Buffer
	binary.Write(&xmitBuf, binary.BigEndian, xmitHdr)
	b = xmitBuf.Bytes()

	return b
}
