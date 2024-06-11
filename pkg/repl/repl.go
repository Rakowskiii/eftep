package repl

import (
	"bufio"
	config "eftep/pkg/config/client"
	"fmt"
	"os"
	"strings"
)

var Socket int
var CurrentConnection string

const (
	Exit          string = "q"
	Connect       string = "conn"
	ConnectManual string = "connm"
	Disconnect    string = "dc"
	Help          string = "?"
	GetFile       string = "get"
	PutFile       string = "put"
	DeleteFile    string = "del"
	ListDir       string = "dir"
	RenameFile    string = "mv"
	Discover      string = "find"
)

func Run() {
	if _, err := os.Stat(config.DOWNLOAD_DIR); os.IsNotExist(err) {
		os.Mkdir(config.DOWNLOAD_DIR, 0755)
	}

	fmt.Println("Eftep Repl -- v0.1")
	for {
		command, err := awaitCommand()
		if err != nil {
			fmt.Println("Failed to read command:", err)
			continue
		}

		handleCommand(command)
	}
}

func awaitCommand() (string, error) {
	prompt()
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", fmt.Errorf("failed to read input: %v", scanner.Err())
	}
	input := strings.TrimSpace(scanner.Text())
	tokens := strings.Split(input, " ")

	return tokens[0], nil
}

func handleCommand(command string) {
	handler, found := commandHandlers[command]
	if found {
		handler()
	} else {
		fmt.Printf("%s is not recognized command. Try %v to get list of available commands.\n", command, Help)
	}
}

var commandHandlers = map[string]func(){
	Connect:       handleConnect,
	ConnectManual: handleConnectManual,
	Disconnect:    handleDisconnect,
	Help:          showHelp,
	GetFile:       func() { handleIfConnected(handleGetFile) },
	PutFile:       func() { handleIfConnected(handleFileUpload) },
	DeleteFile:    func() { handleIfConnected(handleDeleteFile) },
	ListDir:       func() { handleIfConnected(handleListDir) },
	RenameFile:    func() { handleIfConnected(handleRenameFile) },
	Exit:          exit,
	Discover:      handleDiscover,
}

func prompt() {
	fmt.Printf("%s> ", CurrentConnection)
}
