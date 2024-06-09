package main

import (
	"eftep/pkg/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var IP_ADDR = [4]byte{0, 0, 0, 0}

func main() {
	// Setup logging to a /var/log/eftep.log file
	if err := setupLogs(); err != nil {
		log.Panicf("Failed to setup logs: %v", err)
	}

	// If the work directory doesn't exist, create it
	if _, err := os.Stat(server.WORKDIR); os.IsNotExist(err) {
		os.Mkdir(server.WORKDIR, 0755)
	}

	// Ignore SIGHUP signals to work in daemon mode
	signal.Ignore(syscall.SIGHUP)

	// Setup the socket and start listening for connections
	socket, err := setupSocket()
	if err != nil {
		log.Fatalf("Failed to setup socket: %v", err)
	}

	// Start accepting connections, and handle them in a separate goroutines
	// No need for a worker pool, as the server is not expected to handle a lot of clients
	for {
		client, addr, err := syscall.Accept(socket)
		if err != nil {
			log.Println("[Error] Failed to accept connection:", err)
		}
		log.Println("[Info] Accepted connection from:", addr)
		go server.HandleClient(client)
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
		Addr: IP_ADDR,
		Port: server.PORT,
	}

	if err = syscall.Bind(socket, &sockaddr); err != nil {
		return -1, err
	}

	if err = syscall.Listen(socket, 5); err != nil {
		return -1, err
	}

	return socket, nil
}

func setupLogs() error {
	logFile, err := os.OpenFile("/var/log/eftep.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	return nil
}
