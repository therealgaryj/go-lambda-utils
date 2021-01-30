package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var currentLogLevel = 2 // default to INFO
var availableLogLevels = map[string]int{"TRACE": 0, "DEBUG": 1, "INFO": 2, "WARN": 3, "ERROR": 4, "FATAL": 5}
var requestId = ""
var sequenceNumber = 0

type logMessage struct {
	SequenceNumber int `json:"sequenceNumber"`
	RequestId string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Level string `json:"level"`
	Message string `json:"message"`
}

// Instantiate will read environment variable LOG_LEVEL, if not there or invalid, defaults to INFO
func Instantiate(currentRequestId string) {
	providedLogLevel := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	requestId = currentRequestId

	if logLevel, ok := availableLogLevels[providedLogLevel]; ok {
		currentLogLevel = logLevel
		printLog("INFO", "Set logLevel to be: " + providedLogLevel, true)
	} else {

		defaultLogLevel := ""
		for k, v := range availableLogLevels {
			if v == currentLogLevel {
				defaultLogLevel = k
				break
			}
		}
		printLog("WARN", "Failed to set the log level via environment variable, logLevel is: " + defaultLogLevel, true)
	}
}

// Trace logs at log level TRACE where currentLogLevel is set to TRACE
func Trace(message string) {
	printLog("TRACE", message, false)
}

// Debug logs at log level DEBUG where currentLogLevel is set to TRACE, DEBUG
func Debug(message string) {
	printLog("DEBUG", message, false)
}

// Info logs at log level INFO where currentLogLevel is set to TRACE, DEBUG, INFO
func Info(message string) {
	printLog("INFO", message, false)
}

// Warn logs at log level WARN where currentLogLevel is set to TRACE, DEBUG, INFO, WARN
func Warn(message string) {
	printLog("WARN", message, false)
}

// Error logs at log level ERROR where currentLogLevel is set to TRACE, DEBUG, INFO, WARN, ERROR
func Error(message string) {
	printLog("ERROR", message, false)
}

// Fatal logs at log level FATAL where currentLogLevel is set to TRACE, DEBUG, INFO, WARN, ERROR, FATAL
func Fatal(message string) {
	printLog("FATAL", message, false)
}

func printLog(level string, message string, overrideLogLevel bool) {

	if overrideLogLevel ||  availableLogLevels[level] >= currentLogLevel {

		level = padLevelName(level)
		timestamp := time.Now().UTC().Format("2006-01-02 15:04:05.000")

		structuredMessage := logMessage{
			RequestId: requestId,
			Timestamp: timestamp,
			Level:     level,
			Message:   message,
			SequenceNumber: sequenceNumber,
		}
		json, err := json.Marshal(structuredMessage)

		if err != nil {
			fmt.Printf("Failed to log message: %s\n", err.Error())
		}

		fmt.Println(string(json))
		sequenceNumber++
	}
}

func padLevelName(level string) string {
	charsInLevel := len(level)
	if charsInLevel < 5 {
		for i := 5 - charsInLevel; i > 0; i-- {
			level += " "
		}
	}

	return level
}