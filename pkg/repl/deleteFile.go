package repl

import (
	"bufio"
	"eftep/pkg/commons"
	"fmt"
	"os"
)

func handleDeleteFile(socket int) {
	// Read the filenames to rename
	fmt.Print("Enter the filename to delete (example.txt): ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Failed to read filename")
		return
	}

	filename := scanner.Bytes()

	// Send the command to the server
	if err := sendCommand(socket, commons.DeleteFile, filename); err != nil {
		fmt.Println("Failed to send command to server:", err)
		return
	}

	// Read the response from the server
	handleResponse(socket)
}
