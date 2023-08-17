package internal

import (
	"Master/constants"
	"log"
	"os"
)

func WriteOnLog(logMessage string) {
	var config constants.Config
	constants.ReadJsonConfig(&config)
	// Open the file in append mod
	file, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Impossible to open log file: %v", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Impossibile to close log file: %v", err)
		}
	}(file)

	// Write into the log
	_, err = file.WriteString(logMessage)
	if err != nil {
		log.Fatalf("Impossible to write on log: %v", err)
		return
	}

}
