package server

import (
	"context"
	"eftep/pkg/commons"
	config "eftep/pkg/config/server"
	"eftep/pkg/log"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
)

func handlePutFile(ctx context.Context, client int, dataLen int) {
	// Read the filename to upload
	filename := make([]byte, dataLen)
	commons.ReadFull(client, filename)

	log.Info(ctx, "put_file", string(filename))

	// Create the file
	file, err := os.Create(config.WORKDIR + "/" + string(filename))
	if err != nil {
		log.Error(ctx, "creating file", err)
		sendMessage(ctx, client, []byte("Failed to create file"))
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
			log.Error(ctx, "reading file contents", err)
			// TODO: what to do here? Should we kill the connection? Should we return an error message?
			// We defnitely should delete the file
			file.Close()
			os.Remove(config.WORKDIR + "/" + string(filename))
			sendMessage(ctx, client, []byte("Failed to upload file"))
			return
		}

		file.Write(buf[:n])
		fileLen -= commons.Min(uint32(n), fileLen)
	}

	// Send a success message to the client
	sendMessage(ctx, client, []byte("File uploaded successfully"))
}
