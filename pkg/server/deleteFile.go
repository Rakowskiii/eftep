package server

import (
	"context"
	"eftep/pkg/commons"
	config "eftep/pkg/config/server"
	log "eftep/pkg/log"
	"os"
)

func handleDeleteFile(ctx context.Context, client int, dataLen int) {
	// Read the filename to delete
	filename := make([]byte, dataLen)
	commons.ReadFull(client, filename)

	// Delete the file
	err := os.Remove(config.WORKDIR + "/" + string(filename))
	if err != nil {
		log.Error(ctx, "deleting file", err)
		sendMessage(ctx, client, []byte("Failed to delete file. Check if the file exists"))
		return
	}

	log.Info(ctx, "delete_file", string(filename))

	// Send a success message to the client
	sendMessage(ctx, client, []byte("File deleted successfully"))
}
