package main

import (
	"context"
	"eftep/pkg/commons"
	"eftep/pkg/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	config "eftep/pkg/config/server"
	log "eftep/pkg/log"
)

func main() {
	ctx := context.WithValue(context.Background(), log.SessionIDKey, "eftep")
	// Setup logging to a /var/log/eftep.log file
	logFile, err := log.SetupLogs()
	if err != nil {
		panic(fmt.Sprintf("Failed to setup logs: %v\n", err))
	}
	defer logFile.Close()

	// If the work directory doesn't exist, create it
	if _, err := os.Stat(config.WORKDIR); os.IsNotExist(err) {
		os.Mkdir(config.WORKDIR, 0755)
	}

	// Ignore SIGHUP signals to work in daemon mode
	signal.Ignore(syscall.SIGHUP)

	// Setup the socket and start listening for connections
	socket, err := setupSocket()
	if err != nil {
		log.Error(ctx, "setup socket", err)
		panic("Failed to setup socket")
	}

	// Start the discovery service
	go server.DiscoveryService()

	// Start accepting connections, and handle them in a separate goroutines
	// No need for a worker pool, as the server is not expected to handle a lot of clients
	for {
		client, addr, err := syscall.Accept(socket)
		if err != nil {
			log.Error(ctx, "socket accept", err)
		}

		parsedAddr := commons.ParseIpAddr(addr)
		sessId := server.RandId(4)
		log.Info(ctx, "accept_new_client", fmt.Sprintf("%s with session: %s", parsedAddr, sessId))

		cctx := context.WithValue(ctx, log.ClientIPKey, parsedAddr)
		cctx = context.WithValue(cctx, log.SessionIDKey, sessId)
		go server.HandleClient(cctx, client)
	}
}

func setupSocket() (int, error) {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return -1, err
	}

	if err = syscall.SetsockoptInt(socket, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return -1, err
	}

	if err = syscall.SetsockoptInt(socket, syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1); err != nil {
		return -1, err
	}

	sockaddr := syscall.SockaddrInet4{
		Addr: config.IP_ADDR,
		Port: config.EFTEP_PORT,
	}

	if err = syscall.Bind(socket, &sockaddr); err != nil {
		return -1, err
	}

	if err = syscall.Listen(socket, 5); err != nil {
		return -1, err
	}

	return socket, nil
}
