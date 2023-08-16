package constants

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	NumMapper     int     `json:"num_mapper"`
	NumReducer    int     `json:"num_reducer"`
	MapperPN      int     `json:"mapper_pn"`
	ReducerPN     int     `json:"reducer_pn"`
	DampingFactor float32 `json:"damping_factor"`
	Epsilon       float32 `json:"epsilon"`
	GraphPath     string  `json:"graph_path"`
	MaxIteration  int     `json:"max_iteration"`
	NumNodes      int     `json:"num_nodes"`
	EdgesToAttach int     `json:"edges_to_attach"`
	Seed          int64   `json:"seed"`
}

func ReadJsonConfig(config *Config) {
	// Find the file with relative path
	filePath := "./constants/configuration.json"

	configFile, err := os.ReadFile(filePath)
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
