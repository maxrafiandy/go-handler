package handler

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger logs information to log/info.log file
func Logger(info interface{}) {
	// check parent type of info interface

	switch info.(type) {
	// if it's a response
	case Response:
		// get data's data type
		data := info.(Response)

		// check whatever the Data is an error
		switch data.Data.(type) {
		// the response contain error
		case error, Error, *Error:
			logError(info)

		// the response is a map
		case map[string]interface{}:
			logInfo(info)

		// it's a pure response
		default:
			logResponse(info)
		}

	// if its an error
	case error, Error, *Error:
		logError(info)

	// otherwise return it as default response
	default:
		logInfo(info)
	}
}

// logError logs information to log/error.log file
func logError(info interface{}) {
	now := time.Now().Format(formatDateYMD)
	filepath := fmt.Sprintf("log/ERROR_%s.log", now)
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println(info)
}

// logInfo logs information to log/error.log file
func logInfo(info interface{}) {
	now := time.Now().Format(formatDateYMD)
	filepath := fmt.Sprintf("log/DEBUG_%s.log", now)
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println(info)
}

// logResponse logs information to log/info.log file
func logResponse(info interface{}) {
	now := time.Now().Format(formatDateYMD)
	filepath := fmt.Sprintf("log/RESPONSE_%s.log", now)
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println(info)
}
