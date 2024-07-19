package internal

import (
	"fmt"
	"os"
	"time"
)

func formatLog(level, message string) string {
	layout := "2006/01/02 - 15:04:05"
	return fmt.Sprintf("[str] %s | %s | %s", time.Now().Format(layout), level, message)
}

func Fatalf(message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	logMessage := formatLog("FATAL", formattedMessage)
	fmt.Println(logMessage)
	os.Exit(1)
}

func Infof(message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	logMessage := formatLog("INFO", formattedMessage)
	fmt.Println(logMessage)
}

func Warnf(message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	logMessage := formatLog("WARN", formattedMessage)
	fmt.Println(logMessage)
}

func Errorf(message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	logMessage := formatLog("ERROR", formattedMessage)
	fmt.Println(logMessage)
}
