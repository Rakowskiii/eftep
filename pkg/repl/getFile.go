package repl

import (
	"bufio"
	"eftep/pkg/commons"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"

	config "eftep/pkg/config/client"
)

func handleGetFile(socket int) {
	// Read the filename to download
	fmt.Print("Enter the filename to download: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Failed to read filename")
		return
	}

	fileName := scanner.Bytes()

	sendCommand(socket, commons.GetFile, fileName)

	// Read the file size
	fileSize := make([]byte, 4)
	_, err := syscall.Read(socket, fileSize)
	if err != nil {
		fmt.Println("Failed to read file size from server:", err)
		return
	}

	size := binary.BigEndian.Uint32(fileSize)
	if size == 0 {
		fmt.Println("File not found on server")
		return
	}

	fmt.Println("Receiving file of length:", size)

	// Create the file
	file, err := os.Create(config.DOWNLOAD_DIR + "/" + string(fileName))
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	// Read the file contents
	buf := make([]byte, 4096)
	for size > 0 {
		n, err := syscall.Read(socket, buf)
		if err != nil {
			fmt.Println("Failed to read file contents from server:", err)
			return
		}

		if n == 0 {
			break
		}

		file.Write(buf[:n])
		size -= commons.Min(uint32(n), size)
	}

	fmt.Println("File download complete")

}
