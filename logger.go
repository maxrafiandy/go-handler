package handler

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	levelERROR string = "ERROR"
	levelINFO  string = "INFO"
)

// Logger logs information to log/info.log file
func Logger(payload interface{}) {
	// check parent type of payload interface

	createLogDirectory()

	switch payload.(type) {
	// if it's a response
	case Response:
		// get data's data type
		data := payload.(Response)

		// check whatever the Data is an error
		switch data.Data.(type) {
		// the response contain error
		case error, Error, *Error:
			writeLog(levelERROR, payload)

		// the response is a map
		// case map[string]interface{}:
		// 	logInfo(payload)

		// it's a pure response
		default:
			writeLog(levelINFO, payload)
		}

	// if its an error
	case error, Error, *Error:
		writeLog(levelERROR, payload)

	// otherwise return it as default response
	default:
		writeLog(levelINFO, payload)
	}
}

// writeLog logs information to log/info.log file
func writeLog(level string, payload interface{}) {
	var (
		err      error
		now      string
		filepath string
		file     *os.File
	)

	if level == "" {
		level = "INFO"
	}

	now = time.Now().Format(formatDateYMD)
	filepath = fmt.Sprintf("log/%s_%s.log", level, now)

	if file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.Fatalf("[go-handler] error opening file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)

	log.Printf("[go-handler] %s: %+v\n", level, payload)
}

func createLogDirectory() {
	var err error
	if _, err = os.Stat("log"); os.IsNotExist(err) {
		if err = os.Mkdir("log", os.ModeDir); err != nil {
			log.Fatalf("[go-handler] fatal: %v", err)
		}
	}
}
