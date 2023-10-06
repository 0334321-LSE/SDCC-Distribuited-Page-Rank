package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	NumMapper  int `json:"num_mapper"`
	NumReducer int `json:"num_reducer"`
	MapperPN   int `json:"mapper_pn"`
	ReducerPN  int `json:"reducer_pn"`
}

func ReadJsonConfig(config *Config) {
	// Find the file with relative path
	filePath := filepath.Join("..", "configuration.json")

	// Obtain abs path
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatal(err)
	}

	configFile, err := os.ReadFile(absFilePath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println("Error parsing config:", err)
		return
	}
}
