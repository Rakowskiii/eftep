package commons

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"time"
)

const DISCOVERY_MESSAGE = "WHERE_ARE_YOU"
const DISCOVERY_RESPONSE = "I_AM_EFTEP"
const DISCOVERY_TIMEOUT_SECS = 1 * time.Second

const (
	ListDir    = 0
	GetFile    = 1
	PutFile    = 2
	DeleteFile = 3
	RenameFile = 4
)

func ReadFull(client int, buf []byte) (int, error) {
	total := 0
	for total < len(buf) {
		n, err := syscall.Read(client, buf[total:])
		if err != nil {
			return total, err
		}
		if n == 0 {
			return total, fmt.Errorf("client disconnected")
		}
		total += n
	}
	return total, nil
}

func MakeMessage(data []byte) []byte {
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(data)))
	return append(lenBytes, data...)
}

func Min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func ParseIpAddr(addr syscall.Sockaddr) string {
	switch addr := addr.(type) {
	case *syscall.SockaddrInet4:
		return parseInetIp4Addr(addr)
	case *syscall.SockaddrInet6:
		return parseInetIp6Addr(addr)
	default:
		return fmt.Sprintf("%s", addr)
	}
}

func parseInetIp4Addr(addr *syscall.SockaddrInet4) string {
	return fmt.Sprintf("%v.%v.%v.%v", addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3])
}

func parseInetIp6Addr(addr *syscall.SockaddrInet6) string {
	return fmt.Sprintf("%v.%v.%v.%v", addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3])
}
