package repl

import (
	"bufio"
	"eftep/pkg/commons"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var Socket int

var DownloadsDir = "/tmp/eftepcli"

func Run() {
	if _, err := os.Stat(DownloadsDir); os.IsNotExist(err) {
		os.Mkdir(DownloadsDir, 0755)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Eftep Repl")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		handleCommand(input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}

const (
	Exit       string = "q"
	Connect    string = "conn"
	Disconnect string = "dc"
	Help       string = "?"
	GetFile    string = "get"
	PutFile    string = "put"
	DeleteFile string = "del"
	ListDir    string = "dir"
	RenameFile string = "mv"
	Discover   string = "find"
)

var commandHandlers = map[string]func(){
	Connect:    handleConnect,
	Disconnect: handleDisconnect,
	Help:       showHelp,
	GetFile:    func() { handleIfConnected(handleGetFile) },
	PutFile:    func() { handleIfConnected(handleFileUpload) },
	DeleteFile: func() { handleIfConnected(handleDeleteFile) },
	ListDir:    func() { handleIfConnected(handleListDir) },
	RenameFile: func() { handleIfConnected(handleRenameFile) },
	Exit:       exit,
	Discover:   func() { fmt.Println("Not implemented") },
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

func handleDisconnect() {
	if Socket == 0 {
		fmt.Println("Not connected to server")
		return
	}
	syscall.Close(Socket)
	Socket = 0
	fmt.Println("Disconnected from server")
}

func handleCommand(line string) {
	tokens := strings.Split(line, " ")
	if len(tokens) == 0 {
		fmt.Println("No command entered")
		return
	}
	fmt.Println("Command:", tokens[0])

	handler, found := commandHandlers[tokens[0]]
	if found {
		handler()
	} else {
		fmt.Println("Unknown command:", tokens[0])
	}
}

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

func handleListDir(socket int) {
	// Send the command to the server
	command := []byte{0, 0, 0, 0, 0}
	if _, err := syscall.Write(socket, command); err != nil {
		fmt.Println("Failed to send command to server:", err)
		return
	}

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

func handleConnect() {
	// Read server address and port from stdin
	fmt.Print("Enter server address (e.g., 127.0.0.1): ")
	var serverAddr string
	n, err := fmt.Scanln(&serverAddr)

	if err != nil || n == 0 {
		fmt.Println("Defaulting to 127.0.0.1")
		serverAddr = "127.0.0.1"
	}

	fmt.Print("Enter server port (e.g., 8080): ")
	var serverPort int
	n, err = fmt.Scanln(&serverPort)
	if err != nil || n == 0 {
		fmt.Println("Defaulting to port 8080")
		serverPort = 8080
	}

	// Convert server address to [4]byte
	ipParts := strings.Split(serverAddr, ".")
	if len(ipParts) != 4 {
		fmt.Println("Invalid IP address format")
		return
	}

	var ipAddr [4]byte
	for i := 0; i < 4; i++ {
		part, err := strconv.Atoi(ipParts[i])
		if err != nil {
			fmt.Println("Invalid IP address part:", err)
			return
		}
		ipAddr[i] = byte(part)
	}

	// Create a socket
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("Failed to create socket:", err)
		return
	}

	// Prepare the sockaddr for connecting
	sockaddr := &syscall.SockaddrInet4{
		Addr: ipAddr,
		Port: serverPort,
	}

	// Connect to the server (binds to the first available port)
	if err := syscall.Connect(socket, sockaddr); err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}

	// Set the global Socket variable
	Socket = socket
	fmt.Println("Connected to server successfully")
}

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

func handleGetFile(socket int) {
	// Read the filename to download
	fmt.Print("Enter the filename to download: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Failed to read filename")
		return
	}

	fileName := scanner.Text()

	message := commons.MakeMessage([]byte(fileName))

	// Send the command to the server
	message = append([]byte{commons.GetFile}, message...)
	if _, err := syscall.Write(socket, message); err != nil {
		fmt.Println("Failed to send command to server:", err)
		return
	}

	// Read the file size
	fileSize := make([]byte, 4)
	_, err := syscall.Read(socket, fileSize)
	if err != nil {
		fmt.Println("Failed to read file size from server:", err)
		return
	}

	size := binary.BigEndian.Uint32(fileSize)
	fmt.Println("Receiving file of length:", size)

	// Create the file
	file, err := os.Create(DownloadsDir + "/" + (fileName))
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
