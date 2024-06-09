package repl

import (
	"bufio"
	config "eftep/pkg/config/client"
	"fmt"
	"os"
	"strings"
)

var Socket int

const (
	Exit       string = "q"
	Connect    string = "conn"
	Disconnect string = "dc"
	Help       string = "?"
	GetFile    string = "get"
	PutFile    string = "put"
	DeleteFile string = "del"
	ListDir    string = "dir"
	RenameFile string = "mv"
	Discover   string = "find"
)

func Run() {
	if _, err := os.Stat(config.DOWNLOAD_DIR); os.IsNotExist(err) {
		os.Mkdir(config.DOWNLOAD_DIR, 0755)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Eftep Repl")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		handleCommand(input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}

func handleCommand(line string) {
	tokens := strings.Split(line, " ")
	if len(tokens) == 0 {
		fmt.Println("No command entered")
		return
	}
	fmt.Println("Command:", tokens[0])

	handler, found := commandHandlers[tokens[0]]
	if found {
		handler()
	} else {
		fmt.Println("Unknown command:", tokens[0])
	}
}

var commandHandlers = map[string]func(){
	Connect:    handleConnect,
	Disconnect: handleDisconnect,
	Help:       showHelp,
	GetFile:    func() { handleIfConnected(handleGetFile) },
	PutFile:    func() { handleIfConnected(handleFileUpload) },
	DeleteFile: func() { handleIfConnected(handleDeleteFile) },
	ListDir:    func() { handleIfConnected(handleListDir) },
	RenameFile: func() { handleIfConnected(handleRenameFile) },
	Exit:       exit,
	Discover:   handleDiscover,
}
