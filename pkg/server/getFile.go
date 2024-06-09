package server

import (
	"context"
	"eftep/pkg/commons"
	config "eftep/pkg/config/server"
	log "eftep/pkg/log"
	"encoding/binary"
	"os"
	"syscall"
)

func handleGetFile(ctx context.Context, socket int, dataLen int) {
	// Read the filename
	filename := make([]byte, dataLen)
	commons.ReadFull(socket, filename)

	log.Info(ctx, "get_file", string(filename))

	// Open the file
	file, err := os.Open(config.WORKDIR + "/" + string(filename))
	if err != nil {
		log.Error(ctx, "opening file", err)
		socketWrite(ctx, socket, []byte{0, 0, 0, 0})
		return
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		log.Error(ctx, "getting file info", err)
		return
	}
	fileSize := uint32(fileInfo.Size())

	// Send the file size to the client
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(fileSize))
	if _, err := syscall.Write(socket, sizeBytes); err != nil {
		log.Error(ctx, "writing file size", err)
		return
	}

	// Send the file contents to the client
	buf := make([]byte, 4096)
	for fileSize > 0 {
		n, err := file.Read(buf)
		if err != nil {
			log.Error(ctx, "reading file contents", err)
			return
		}

		if n == 0 {
			break
		}

		if _, err := syscall.Write(socket, buf[:n]); err != nil {
			log.Error(ctx, "writing file contents", err)
			return
		}

		fileSize -= commons.Min(uint32(n), uint32(fileSize))
	}
}
