package repl

import (
	config "eftep/pkg/config/client"
	"fmt"
	"strconv"
	"strings"
	"syscall"
)

func connect(ip [4]byte, port int) {

	if err := setupSocket(ip, port); err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}

	fmt.Println("Connected to server successfully")
}

func handleConnect() {
	if len(KnownHosts) == 0 {
		fmt.Println("No known hosts to connect to")
		return
	}

	list := make([]string, 0, len(KnownHosts))
	for host := range KnownHosts {
		list = append(list, host)
	}

	for i, host := range list {
		fmt.Printf("[%v] %s\n", i, host)
	}

	idx, err := readTargetIndex()
	if err != nil {
		fmt.Println("Failed to read target index:", err)
		return
	}

	if idx < 0 || idx >= len(list) {
		fmt.Println("Invalid target index")
		return
	}

	ip, port, err := hostnameToIpPort(list[idx])
	if err != nil {
		fmt.Println("Failed to parse hostname:", err)
		return
	}

	connect(ip, port)
}

func hostnameToIpPort(hostname string) ([4]byte, int, error) {
	parts := strings.Split(hostname, "@")
	parts = strings.Split(parts[1], ":")
	if len(parts) != 2 {
		return [4]byte{}, 0, fmt.Errorf("invalid hostname format")
	}

	ip, err := parseIp(parts[0])
	if err != nil {
		return [4]byte{}, 0, fmt.Errorf("invalid IP address: %v", err)
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return [4]byte{}, 0, fmt.Errorf("invalid port: %v", err)
	}

	return ip, port, nil
}

func handleConnectManual() {
	ip, err := readIpAddr()
	if err != nil {
		fmt.Println("Failed to read server address:", err)
		return
	}

	port, err := readPort()
	if err != nil {
		fmt.Println("Failed to read server port:", err)
		return
	}

	connect(ip, port)
}

func handleDisconnect() {
	if Socket == 0 {
		fmt.Println("Not connected to server")
		return
	}
	syscall.Close(Socket)
	Socket = 0
	CurrentConnection = ""
	fmt.Println("Disconnected from server")
}

func setupSocket(ip [4]byte, port int) error {
	// Create a socket
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("failed to create socket: %v", err)
	}

	// Set the receive timeout to 1 second
	tv := syscall.NsecToTimeval(config.RECV_TIMEOUT_SECS.Nanoseconds())
	if err = syscall.SetsockoptTimeval(socket, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv); err != nil {
		return fmt.Errorf("failed to set SO_RCVTIMEO: %v", err)
	}

	// Prepare the sockaddr for connecting
	sockaddr := &syscall.SockaddrInet4{
		Addr: ip,
		Port: port,
	}

	// Connect to the server (binds to the first available port)
	if err := syscall.Connect(socket, sockaddr); err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	Socket = socket
	CurrentConnection = fmt.Sprintf("[%v.%v.%v.%v:%v]", ip[0], ip[1], ip[2], ip[3], port)

	return nil
}

func readPort() (int, error) {
	fmt.Print("Enter server port (e.g., 8080): ")
	var serverPort int
	n, err := fmt.Scanln(&serverPort)
	if err != nil || n == 0 {
		fmt.Println("Defaulting to port 8080")
		return 8080, nil
	}

	return serverPort, nil
}

func readTargetIndex() (int, error) {
	var targetIndex int
	fmt.Print("Enter the index of the server to connect to: ")
	n, err := fmt.Scanln(&targetIndex)
	if err != nil || n == 0 {
		return 0, fmt.Errorf("failed to read target index")
	}

	return targetIndex, nil
}

func readIpAddr() ([4]byte, error) {
	// Read server address and port from stdin
	fmt.Print("Enter server address (e.g., 127.0.0.1): ")
	var serverAddr string
	n, err := fmt.Scanln(&serverAddr)
	if err != nil || n == 0 {
		fmt.Println("Defaulting to 127.0.0.1")
		return [4]byte{127, 0, 0, 1}, nil
	}

	return parseIp(serverAddr)
}

func parseIp(serverAddr string) ([4]byte, error) {
	ipParts := strings.Split(serverAddr, ".")
	if len(ipParts) != 4 {
		return [4]byte{}, fmt.Errorf("invalid IP address format")
	}

	var ipAddr [4]byte
	for i := 0; i < 4; i++ {
		part, err := strconv.Atoi(ipParts[i])
		if err != nil {
			return [4]byte{}, fmt.Errorf("invalid IP address part: %v", err)
		}
		ipAddr[i] = byte(part)
	}

	return ipAddr, nil
}
