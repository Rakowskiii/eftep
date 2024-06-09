package server

import (
	"context"
	"eftep/pkg/commons"
	"eftep/pkg/log"
	"errors"
	"fmt"
	"syscall"
)

func HandleClient(ctx context.Context, client int) {
	defer syscall.Close(client)

	for {
		// Read the 5-byte header
		command, dataLen, err := readHeader(client)
		if err != nil {
			log.Error(ctx, "reading header", err)
			return
		}

		// Pass the data to the appropriate handler based on the action
		switch command {
		case commons.ListDir:
			handleListDir(ctx, client)
		case commons.GetFile:
			handleGetFile(ctx, client, int(dataLen))
		case commons.PutFile:
			handlePutFile(ctx, client, int(dataLen))
		case commons.DeleteFile:
			handleDeleteFile(ctx, client, int(dataLen))
		case commons.RenameFile:
			handleRenameFile(ctx, client, int(dataLen))
		default:
			log.Error(ctx, "handle command", errors.New(fmt.Sprintf("unknown action: %d", command)))
		}
	}
}
