package server

import (
	"context"
	config "eftep/pkg/config/server"
	"eftep/pkg/log"
	"os"
)

func handleListDir(ctx context.Context, socket int) {
	// List files in the work directory
	files, err := os.ReadDir(config.WORKDIR)
	if err != nil {
		log.Error(ctx, "reading directory", err)
		sendMessage(ctx, socket, []byte("Failed to list files"))
		return
	}

	log.Info(ctx, "list_files", config.WORKDIR)

	// Send the list of filenames to the socket
	var filenames []byte
	for _, file := range files {
		filenames = append(filenames, []byte(file.Name()+" ")...)
		filenames = append(filenames, 0)
	}

	// Send the length of the filenames list
	sendMessage(ctx, socket, filenames)
}
