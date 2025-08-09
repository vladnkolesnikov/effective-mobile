package logger

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

func LogError(msg string, err error) {
	logger.Printf("[ERROR] %s: %v", msg, err)
}

func LogInfo(msg string) {
	logger.Printf("[INFO] %s", msg)
}

func LogFatal(err error) {
	logger.Fatalf("[FATAL] %v", err)
}
