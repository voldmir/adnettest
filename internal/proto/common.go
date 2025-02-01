package proto

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/voldmir/adnettest/internal/common"
)

type Proto interface {
	Message(remoteAddr net.IP, localAddr net.IP, localPort int) []byte
	Protocols() []string
	Port() string
	Buff() []byte
}

func TestService(host string, proto Proto, timeout int64) {
	var err error
	var conn net.Conn
	for _, protocol := range proto.Protocols() {

		port := proto.Port()
		_timeout := time.Second * time.Duration(timeout)

		conn, err = net.DialTimeout(protocol, net.JoinHostPort(host, port), _timeout)
		if err != nil {
			goto FAIL
		}
		defer conn.Close()

		if protocol == "udp" {

			buff := proto.Buff()

			msg := proto.Message(conn.RemoteAddr().(*net.UDPAddr).IP.To4(),
				conn.LocalAddr().(*net.UDPAddr).IP.To4(),
				conn.LocalAddr().(*net.UDPAddr).Port)

			conn.SetDeadline(time.Now().Add(_timeout))

			_, err = conn.Write(msg)
			if err != nil {
				goto FAIL
			} else {
				_, err = bufio.NewReader(conn).Read(buff)
				if err != nil {
					goto FAIL
				}
			}
		}

		if conn != nil {
			str := common.IndentByRuneRight(fmt.Sprintf("%s/%s", port, protocol), 20, ' ')
			fmt.Printf("%sopen\n", str)
			continue
		}
	FAIL:
		str := common.IndentByRuneRight(fmt.Sprintf("%s/%s", port, protocol), 20, ' ')
		fmt.Printf("%sclosed\t%+v\n", str, err)
		continue
	}
}
