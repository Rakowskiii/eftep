package repl

import (
	"bufio"
	"eftep/pkg/commons"
	"fmt"
	"os"
	"strings"
	"syscall"
)

func handleRenameFile(socket int) {
	// Read the filenames to rename
	fmt.Print("Enter the space separated filenames to rename (oldname newname): ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Failed to read filenames")
		return
	}
	names := strings.Replace(scanner.Text(), " ", ":", 1)

	// convert names to bytes
	namesBytes := []byte(names)
	message := commons.MakeMessage(namesBytes)

	// Send the command to the server
	message = append([]byte{commons.RenameFile}, message...)
	if _, err := syscall.Write(socket, message); err != nil {
		fmt.Println("Failed to send command to server:", err)
		return
	}

	// Read the response from the server
	response := make([]byte, 4096)
	n, err := syscall.Read(socket, response)
	if err != nil {
		fmt.Println("Failed to read response from server:", err)
		return
	}

	fmt.Println("Response from server:", string(response[:n]))
}
