package repl

import (
	config "eftep/pkg/config/client"
	"fmt"
	"strconv"
	"strings"
	"syscall"
)

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

	// Set the receive timeout to 1 second
	tv := syscall.NsecToTimeval(config.RECV_TIMEOUT_SECS.Nanoseconds())
	if err = syscall.SetsockoptTimeval(socket, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv); err != nil {
		fmt.Println("Failed to set SO_RCVTIMEO:", err)
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

func handleDisconnect() {
	if Socket == 0 {
		fmt.Println("Not connected to server")
		return
	}
	syscall.Close(Socket)
	Socket = 0
	fmt.Println("Disconnected from server")
}
