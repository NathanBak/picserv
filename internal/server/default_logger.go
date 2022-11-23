package server

import (
	"fmt"
)

type defaultLogger struct{}

type messageType int8

const (
	_ messageType = iota
	debugMessage
	infoMessage
	warningMessage
	errorMessage
)

func (mt messageType) String() string {
	switch mt {
	case debugMessage:
		return "Debug"
	case infoMessage:
		return "Info"
	case warningMessage:
		return "Warning"
	case errorMessage:
		return "Error"
	}
	return "Unknown"
}

func (log defaultLogger) Debug(msg string) {
	log.log(debugMessage, msg)
}

func (log defaultLogger) Info(msg string) {
	log.log(infoMessage, msg)
}

func (log defaultLogger) Warning(msg string) {
	log.log(warningMessage, msg)
}

func (log defaultLogger) Error(msg string) {
	log.log(errorMessage, msg)
}

func (log defaultLogger) log(msgType messageType, msg string) {
	fmt.Printf("%s: %s\n", msgType.String(), msg)
}
