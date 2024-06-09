package server

import (
	"context"
	"eftep/pkg/commons"
	"eftep/pkg/log"
	"encoding/binary"
	"fmt"
	"math/rand"
	"syscall"
	"time"
)

func socketWrite(ctx context.Context, socket int, message []byte) {
	if _, err := syscall.Write(socket, message); err != nil {
		// Close the client connection if we cant write to it
		syscall.Close(socket)
		log.Error(ctx, "writing message", err)
	}
}

func readHeader(client int) (byte, uint32, error) {
	// Read the 5-byte header
	header := make([]byte, 5)
	n, err := syscall.Read(client, header)
	if err != nil {
		return 0, 0, err
	}
	if n == 0 {
		return 0, 0, fmt.Errorf("client closed connection")
	}

	// Parse the header
	command := header[0]
	dataLen := binary.BigEndian.Uint32(header[1:5])

	return command, dataLen, nil
}

// sendMessage sends a message to the client
func sendMessage(ctx context.Context, socket int, message []byte) {
	message = commons.MakeMessage(message)
	socketWrite(ctx, socket, message)
}

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandId(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
