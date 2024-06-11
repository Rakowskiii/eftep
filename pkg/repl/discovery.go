package repl

import (
	"eftep/pkg/commons"
	config "eftep/pkg/config/client"
	"fmt"
	"strings"
	"syscall"
)

var KnownHosts = make(map[string]bool)

const (
	DiscoveryPart = 0
	NamePart      = 1
	PortPart      = 2
)

func handleDiscover() {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		fmt.Println("Failed to create socket:", err)
	}
	defer syscall.Close(socket)

	if err = syscall.SetsockoptInt(socket, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		fmt.Println("Failed to set SO_REUSEADDR:", err)
		return
	}

	if err = syscall.SetsockoptInt(socket, syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1); err != nil {
		fmt.Println("Failed to set SO_REUSEPORT:", err)
		return
	}

	timeout := commons.DISCOVERY_TIMEOUT_SECS
	tv := syscall.NsecToTimeval(timeout.Nanoseconds())
	if err = syscall.SetsockoptTimeval(socket, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv); err != nil {
		fmt.Println("Failed to set SO_RCVTIMEO:", err)
		return
	}

	sockaddr := syscall.SockaddrInet4{
		Addr: config.DISCOVERY_BIND_IP_ADDR,
		Port: config.DISCOVERY_BIND_PORT,
	}

	if err = syscall.Bind(socket, &sockaddr); err != nil {
		fmt.Println("Failed to bind socket:", err)
		return
	}

	sendMulticastDiscovery(socket, config.MULTICAST_GROUPS[:])
	fmt.Println("Sent multicast discovery message")

	// Reset the known hosts
	KnownHosts = make(map[string]bool)

	// Listen for responses
	for {
		buf := make([]byte, 4096)
		n, addr, err := syscall.Recvfrom(socket, buf, 0)
		if err != nil {
			fmt.Println("Finished waiting for responses")
			break
		}

		// Parse the response (CONFIRMATION:NAME:PORT) to confirm it is a discovery response, and get the server port
		response := strings.Split(string(buf[:n]), ":")

		if len(response) != 3 {
			// Ignore messages that are not valid discovery responses
			continue
		}

		if response[DiscoveryPart] != commons.DISCOVERY_RESPONSE {
			// Ignore messages that are not discovery responses
			continue
		}

		host := fmt.Sprintf("%s@%s:%s", response[NamePart], commons.ParseIpAddr(addr), response[PortPart])
		if _, alreadyKnow := KnownHosts[host]; !alreadyKnow {
			fmt.Println(" Discovered host:", host)
			KnownHosts[host] = true
		}
	}

}

func sendMulticastDiscovery(socket int, addrs [][4]byte) {
	message := []byte(commons.DISCOVERY_MESSAGE)

	for _, addr := range addrs {
		addr := syscall.SockaddrInet4{
			Port: config.DISCOVERY_SERVER_PORT,
			Addr: addr,
		}

		err := syscall.Sendto(socket, message, 0, &addr)
		if err != nil {
			fmt.Printf("Failed to send message: %v for group: %v\n", err, addr)
		}
	}
}
