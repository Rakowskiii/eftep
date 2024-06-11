package repl

import (
	"eftep/pkg/commons"
	"encoding/binary"
	"fmt"
	"syscall"
)

func handleListDir(socket int) {
	sendCommand(socket, commons.ListDir, nil)

	// Read the response from the server
	responseSize := make([]byte, 4)
	_, err := syscall.Read(socket, responseSize)
	if err != nil {
		fmt.Println("Failed to read response from server:", err)
		return
	}

	size := binary.BigEndian.Uint32(responseSize)
	response := make([]byte, size)
	commons.ReadFull(socket, response)

	fmt.Println("Response from server:", string(response[:size]))
}
