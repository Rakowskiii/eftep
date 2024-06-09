package log

import (
	"context"
	"fmt"
	"log"
	"os"

	config "eftep/pkg/config/server"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorArg    = "\033[95m"
)

type contextKey string

const (
	SessionIDKey contextKey = "sessionID"
	ClientIPKey  contextKey = "clientIP"
)

func SetupLogs() (*os.File, error) {
	logFile, err := os.OpenFile(config.LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(logFile)

	Info(context.Background(), "setup logs", fmt.Sprintf("into file %s", config.LOGFILE))

	return logFile, nil
}

func formatMessage(ctx context.Context, src, color, level, msg string) string {
	connectionId, CIDExists := ctx.Value(SessionIDKey).(string)
	clientIp, CIPExists := ctx.Value(ClientIPKey).(string)
	client := fmt.Sprintf(" %s[local@localhost]%s ", ColorYellow, ColorReset)

	if CIDExists {
		if CIPExists {
			client = fmt.Sprintf(" %s[%s@%s]%s ", ColorGreen, connectionId, clientIp, ColorReset)

		} else {
			client = fmt.Sprintf(" %s[%s@localhost]%s ", ColorYellow, connectionId, ColorReset)
		}
	}
	level = fmt.Sprintf("[%s]", level)
	return fmt.Sprintf("%s%7s%s%s%24s%s%s %s%s: %s%s%s\n", color, level, ColorReset, ColorYellow, client, ColorReset, ColorBlue, src, ColorReset, ColorCyan, msg, ColorReset)
}

// Error logs an error with a given context message in red
func Error(ctx context.Context, src string, err error) {
	log.Printf(formatMessage(ctx, src, ColorRed, "Error", err.Error()))
}

// Info logs an info message with a given context in blue
func Info(ctx context.Context, src, msg string) {
	log.Printf(formatMessage(ctx, src, ColorBlue, "Info ", msg))
}

// Debug logs a debug message with a given context in cyan
func Debug(ctx context.Context, src, msg string) {
	log.Printf(formatMessage(ctx, src, ColorCyan, "Debug", msg))
}

// Fatal logs a fatal error with a given context message in purple and exits the program
func Fatal(ctx context.Context, src string, err error) {
	log.Fatalf(formatMessage(ctx, src, ColorPurple, "Fatal", err.Error()))
}
