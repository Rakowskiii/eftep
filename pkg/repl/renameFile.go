package repl

import (
	"bufio"
	"eftep/pkg/commons"
	"fmt"
	"os"
	"strings"
)

func handleRenameFile(socket int) {
	// Read the filenames to rename
	fmt.Print("Enter the space separated filenames to rename (oldname newname): ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Failed to read filenames")
		return
	}
	// TODO: Verify correct input
	names := strings.Replace(scanner.Text(), " ", ":", 1)

	if err := sendCommand(socket, commons.RenameFile, []byte(names)); err != nil {
		fmt.Println("Failed to send command to server:", err)
		return
	}

	// Read the response from the server
	handleResponse(socket)
}
