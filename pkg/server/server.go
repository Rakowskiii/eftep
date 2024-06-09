package server

import (
	"eftep/pkg/commons"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"syscall"
)

const PORT = 8080
const WORKDIR = "/tmp/eftep"

func HandleClient(client int) {
	defer syscall.Close(client)

	for {
		// Read the 5-byte header
		command, dataLen, err := readHeader(client)
		if err != nil {
			fmt.Println("Error reading header:", err)
			return
		}

		// Pass the data to the appropriate handler based on the action
		switch command {
		case commons.ListDir:
			handleListDir(client)
		case commons.GetFile:
			handleGetFile(client, int(dataLen))
		case commons.PutFile:
			handlePutFile(client, int(dataLen))
		case commons.DeleteFile:
			handleDeleteFile(client, int(dataLen))
		case commons.RenameFile:
			handleRenameFile(client, int(dataLen))
		default:
			fmt.Println("Unknown action:", command)
		}
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
		return 0, 0, fmt.Errorf("Client closed connection")
	}

	// Parse the header
	command := header[0]
	dataLen := binary.BigEndian.Uint32(header[1:5])

	return command, dataLen, nil
}

func handleListDir(client int) {
	// List files in the work directory
	files, err := os.ReadDir(WORKDIR)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Send the list of filenames to the client prefixed by the amount of bytes in the list
	var filenames []byte
	for _, file := range files {
		fmt.Println("File:", file.Name())
		filenames = append(filenames, []byte(file.Name()+" ")...)
		filenames = append(filenames, 0)
	}

	// Send the length of the filenames list
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(filenames)))

	message := append(lenBytes, filenames...)
	if _, err := syscall.Write(client, message); err != nil {
		fmt.Println("Error writing length:", err)
		return
	}
}

func handleRenameFile(client int, dataLen int) {
	// Read the filename to rename in loop until all data is read
	filenames := make([]byte, dataLen)
	commons.ReadFull(client, filenames)
	names := strings.Split(string(filenames), ":")

	// Rename the file
	err := os.Rename(WORKDIR+"/"+names[0], WORKDIR+"/"+names[1])
	if err != nil {
		fmt.Println("Error renaming file:", err)
		return
	}

	// Send a success message to the client
	message := []byte("File renamed successfully")
	message = commons.MakeMessage(message)

	if _, err := syscall.Write(client, message); err != nil {
		fmt.Println("Error writing success message:", err)
		return
	}
}
func handleDeleteFile(client int, dataLen int) {
	// Read the filename to rename in loop until all data is read
	filename := make([]byte, dataLen)
	commons.ReadFull(client, filename)

	// Rename the file
	err := os.Remove(WORKDIR + "/" + string(filename))
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}

	// Send a success message to the client
	message := []byte("File deleted successfully")
	message = commons.MakeMessage(message)

	if _, err := syscall.Write(client, message); err != nil {
		fmt.Println("Error writing success message:", err)
		return
	}
}

func handlePutFile(client int, dataLen int) {
	// Read the filename to rename in loop until all data is read
	filename := make([]byte, dataLen)
	commons.ReadFull(client, filename)

	// Create the file
	file, err := os.Create(WORKDIR + "/" + string(filename))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Read next 4 bytes to get the length of the file
	lenBytes := make([]byte, 4)
	commons.ReadFull(client, lenBytes)
	fileLen := binary.BigEndian.Uint32(lenBytes)

	// Read the file contents
	fmt.Println("Receiving file of length:", fileLen)
	buf := make([]byte, 4096)
	for fileLen > 0 {
		n, err := syscall.Read(client, buf)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		file.Write(buf[:n])
		fileLen -= commons.Min(uint32(n), fileLen)
	}

	// Send a success message to the client
	message := []byte("File uploaded successfully")
	message = commons.MakeMessage(message)

	if _, err := syscall.Write(client, message); err != nil {
		fmt.Println("Error writing success message:", err)
		return
	}
}

func handleGetFile(socket int, dataLen int) {
	// Open the file
	filename := make([]byte, dataLen)
	commons.ReadFull(socket, filename)
	file, err := os.Open(WORKDIR + "/" + string(filename))
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	fileSize := uint32(fileInfo.Size())

	// Send the file size to the client
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(fileSize))
	if _, err := syscall.Write(socket, sizeBytes); err != nil {
		fmt.Println("Error writing file size:", err)
		return
	}

	// Send the file contents to the client
	buf := make([]byte, 4096)
	for fileSize > 0 {
		n, err := file.Read(buf)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		if n == 0 {
			break
		}

		if _, err := syscall.Write(socket, buf[:n]); err != nil {
			fmt.Println("Error writing file contents:", err)
			return
		}

		fileSize -= commons.Min(uint32(n), uint32(fileSize))
	}

}
