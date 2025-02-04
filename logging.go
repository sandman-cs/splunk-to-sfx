package main

import (
	"log"
)

// logDebug logs debug messages
func logDebug(message ...any) {
	if conf.LogLevel == "debug" {
		log.Println("DEBUG:", message)
	}
}

// logInfo logs info messages
func logInfo(message ...any) {
	if conf.LogLevel == "debug" || conf.LogLevel == "info" {
		log.Println("INFO:", message)
	}
}

// logError logs error messages
func logError(message ...any) {
	if conf.LogLevel == "debug" || conf.LogLevel == "info" || conf.LogLevel == "error" {
		log.Println("ERROR:", message)
	}
}

// logFatal logs fatal messages
func logFatal(message string, err error) {
	if err != nil {
		log.Fatal("FATAL:", message, err)
	}

}

// logWarn logs warning messages
func logWarn(message ...any) {
	if conf.LogLevel == "info" || conf.LogLevel == "error" || conf.LogLevel == "fatal" || conf.LogLevel == "warn" || conf.LogLevel == "debug" {
		log.Println("WARN:", message)
	}
}

// logMessage logs messages
func logMessage(message ...any) {
	log.Println("MESSAGE:", message)
}
