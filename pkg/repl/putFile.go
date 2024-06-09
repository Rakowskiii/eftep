package repl

import (
	"bufio"
	"eftep/pkg/commons"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

func handleFileUpload(socket int) {
	fmt.Printf("Enter the filename to upload: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Failed to read filename")
		return
	}
	filename := scanner.Text()
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	// Read the file contents
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Failed to get file info:", err)
		return
	}

	message := commons.MakeMessage([]byte(filepath.Base(file.Name())))

	// Send the command to the server
	message = append([]byte{commons.PutFile}, message...)
	if _, err := syscall.Write(socket, message); err != nil {
		fmt.Println("Failed to send command to server:", err)
		return
	}

	fileSize := make([]byte, 4)
	size := fileInfo.Size()
	// encode the file size as a 4-byte big endian integer
	binary.BigEndian.PutUint32(fileSize, uint32(size))
	if _, err := syscall.Write(socket, fileSize); err != nil {
		fmt.Println("Failed to send file size to server:", err)
		return
	}

	// Loop over file contents and send to server in 4kb chunks
	buf := make([]byte, 4096)
	sent := 0
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			fmt.Println("Finished reading file")
			break
		}
		if err != nil {
			fmt.Println("Failed to read file:", err)
			return
		}

		if n == 0 {
			break
		}

		x, err := syscall.Write(socket, buf[:n])
		if err != nil {
			fmt.Println("Failed to send file contents to server:", err)
			return
		}

		sent += x
		fmt.Printf("Sending (%s): [%d/%d];\n", filename, sent, size)
	}
	fmt.Println("File upload complete, waiting for server response")

	// Read the response from the server
	response := make([]byte, 4096)
	n, err := syscall.Read(socket, response)
	if err != nil {
		fmt.Println("Failed to read response from server:", err)
		return
	}

	fmt.Println("Response from server:", string(response[:n]))
}
