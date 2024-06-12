package server

import (
	"context"
	"eftep/pkg/commons"
	config "eftep/pkg/config/server"
	"eftep/pkg/log"
	"fmt"
	"os"
	"strings"
)

func handleRenameFile(ctx context.Context, client int, dataLen int) {
	// Read the filenames to rename
	filenames := make([]byte, dataLen)
	commons.ReadFull(client, filenames)

	names := strings.Split(string(filenames), ":")

	// Rename the file
	err := os.Rename(config.WORKDIR+"/"+names[0], config.WORKDIR+"/"+names[1])
	if err != nil {
		log.Error(ctx, "renaming file", err)
		sendMessage(ctx, client, []byte("Failed to rename file. Check if the file exists"))
		return
	}

	log.Info(ctx, "rename_file", fmt.Sprintf("%s -> %s", names[0], names[1]))

	// Send a success message to the client
	sendMessage(ctx, client, []byte("File renamed successfully"))
}
