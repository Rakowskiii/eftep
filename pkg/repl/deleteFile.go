package repl

import (
	"bufio"
	"eftep/pkg/commons"
	"fmt"
	"os"
	"syscall"
)

func handleDeleteFile(socket int) {
	// Read the filenames to rename
	fmt.Print("Enter the filename to delete (example.txt): ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Failed to read filename")
		return
	}

	message := commons.MakeMessage(scanner.Bytes())

	// Send the command to the server
	message = append([]byte{commons.DeleteFile}, message...)
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
