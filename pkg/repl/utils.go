package repl

import (
	"eftep/pkg/commons"
	"fmt"
	"os"
	"syscall"
)

func sendCommand(socket int, command byte, data []byte) error {
	message := commons.MakeMessage(data)
	message = append([]byte{command}, message...)

	_, err := syscall.Write(socket, message)
	return err
}

func awaitResponse(socket int) (string, error) {
	response := make([]byte, 4096)
	n, err := syscall.Read(socket, response)
	return string(response[:n]), err
}

func handleResponse(socket int) {
	response, err := awaitResponse(socket)
	if err != nil {
		fmt.Println("Failed to read response from server:", err)
		return
	}
	fmt.Println("Response from server:", response)
}

func handleIfConnected(handler func(int)) {
	if Socket == 0 {
		fmt.Println("Not connected to server")
		return
	}
	handler(Socket)
}

func exit() {
	fmt.Print("Exiting...")
	if Socket != 0 {
		syscall.Close(Socket)
	}
	os.Exit(0)
}

func showHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  conn  - Connect to a server")
	fmt.Println("  dir   - List files on the server")
	fmt.Println("  mv    - Rename a file on the server")
	fmt.Println("  del   - Delete a file on the server")
	fmt.Println("  put   - Upload a file to the server")
	fmt.Println("  get   - Download a file from the server")
	fmt.Println("  find  - Discover available servers")
	fmt.Println("  ?     - Show this help message")
	fmt.Println("  q     - Quit the client")
}
