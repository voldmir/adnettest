package common

import "strings"

const (
	PROTO_UDP = "udp"
	PROTO_TCP = "tcp"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func IndentByRuneRight(input string, indent int, r byte) string {
	var i int
	if i = indent - len(input); i < 0 {
		i = 0
	}
	return (input + strings.Repeat(string(r), i))[:indent]
}
