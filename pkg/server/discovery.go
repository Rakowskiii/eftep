package server

import (
	"context"
	"eftep/pkg/commons"
	"fmt"
	"syscall"

	config "eftep/pkg/config/server"
	log "eftep/pkg/log"
)

var ServiceName string

func DiscoveryService(name string) {
	ServiceName = name
	ctx := context.WithValue(context.Background(), log.SessionIDKey, "discovery")
	log.Info(ctx, "discovery", fmt.Sprintf("starting the service with name %s", ServiceName))

	socket := setupSocket(ctx)

	defer syscall.Close(socket)

	// Listen for discovery messages and handle them
	buf := make([]byte, 4096)
	for {
		handleDiscoveryMessage(ctx, socket, buf)
	}
}

func joinMulticastGroups(ctx context.Context, socket int) {
	for _, addr := range config.MULTICAST_GROUPS {
		ipMreq := &syscall.IPMreq{
			Multiaddr: addr,
			Interface: config.IP_ADDR,
		}

		if err := syscall.SetsockoptIPMreq(socket, syscall.IPPROTO_IP, syscall.IP_ADD_MEMBERSHIP, ipMreq); err != nil {
			log.Error(ctx, "join_multicast", err)
		}

		log.Info(ctx, "join_multicast", fmt.Sprintf("group: %v.%v.%v.%v", addr[0], addr[1], addr[2], addr[3]))
	}
}

func setupSocket(ctx context.Context) int {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		log.Error(ctx, "create_socket", err)
	}

	joinMulticastGroups(ctx, socket)

	sockaddr := syscall.SockaddrInet4{
		Port: config.DISCOVERY_PORT,
		Addr: config.IP_ADDR,
	}

	if err := syscall.Bind(socket, &sockaddr); err != nil {
		log.Error(ctx, "bind_socket", err)
	}

	return socket
}

func handleDiscoveryMessage(ctx context.Context, socket int, buf []byte) {
	// Wait for discovery message
	n, addr, err := syscall.Recvfrom(socket, buf, 0)
	if err != nil {
		log.Error(ctx, "udp_recv", err)
		return
	}

	// Ignore messages that are not discovery messages
	if string(buf[:n]) != commons.DISCOVERY_MESSAGE {
		return
	}

	// Respond to the discovery message
	message := commons.DISCOVERY_RESPONSE + ":" + ServiceName + ":" + fmt.Sprintf("%d", config.EFTEP_PORT)
	err = syscall.Sendto(socket, []byte(message), 0, addr)
	if err != nil {
		log.Error(ctx, "udp_response", err)
	}

	log.Info(ctx, "discovery_response", fmt.Sprintf("to: %s", commons.ParseIpAddr(addr)))
}
