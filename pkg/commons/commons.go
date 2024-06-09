package commons

import (
	"encoding/binary"
	"fmt"
	"syscall"
)

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
