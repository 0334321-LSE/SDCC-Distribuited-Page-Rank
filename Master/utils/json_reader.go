package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	NumMapper  int `json:"num_mapper"`
	NumReducer int `json:"num_reducer"`
}

func ReadJsonConfig(config Config) {
	configFile, err := os.ReadFile("./../configuration.json")
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
