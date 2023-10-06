package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type Config struct {
	NumMapper     int     `json:"num_mapper"`
	NumReducer    int     `json:"num_reducer"`
	MapperPn      int     `json:"mapper_pn"`
	MapperHbPn    int     `json:"mapper_hb_pn"`
	ReducerPn     int     `json:"reducer_pn"`
	ReducerHbPn   int     `json:"reducer_hb_pn"`
	DampingFactor float32 `json:"damping_factor"`
	Epsilon       float32 `json:"epsilon"`
	GraphPath     string  `json:"graph_path"`
	LogPath       string  `json:"log_path"`
	MaxIteration  int     `json:"max_iteration"`
	NumNodes      int     `json:"num_nodes"`
	EdgesToAttach int     `json:"edges_to_attach"`
	Seed          int64   `json:"seed"`
	Bucket        string  `json:"bucket"`
	Region        string  `json:"region"`
	SaveOnBucket  bool    `json:"save_on_bucket"`
	PrintGraph    bool    `json:"print_graph"`
	Locally       bool    `json:"locally"`
}

func generateDockerCompose(config Config) error {
	var composeTemplate = `
version: '3'
services:
{{ .MapperServices }}{{ .ReducerServices }}
networks:
  my-network:
`
	if config.Locally == false {
		composeTemplate = `
version: '3'
services:
   app-master:
    build:
      context: ./Master
    volumes:
      - ./output:/Master/output
      - ./configuration.json:/Master/constants/configuration.json
    environment:
      - TZ=Europe/Rome
    networks:
      - my-network{{ .MapperServices }}{{ .ReducerServices }}
networks:
  my-network:
`
	}
	tmpl, err := template.New("docker-compose").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}).Parse(composeTemplate)
	if err != nil {
		return err
	}

	mapperServices := ""
	for i := 1; i <= config.NumMapper; i++ {
		mapperServices += fmt.Sprintf(`
   app-mapper-%d:
    build:
      context: ./Mapper
    ports:
      - "%d:%d"  # Assegna porte univoche ai Mapper
      - "%d:%d"  # Porta univoca per servizio di HeartBeat
    environment:
      - TZ=Europe/Rome
      - MAPPER_PORT=%d
      - HB_PORT=%d
    networks:
      - my-network
`, i, config.MapperPn+i, config.MapperPn+i, config.MapperHbPn+i, config.MapperHbPn+i, config.MapperPn+i, config.MapperHbPn+i)
	}

	reducerServices := ""
	for i := 1; i <= config.NumReducer; i++ {
		reducerServices += fmt.Sprintf(`
   app-reducer-%d:
    build:
      context: ./Reducer
    volumes:
      - ./configuration.json:/Reducer/constants/configuration.json
    ports:
      - "%d:%d"  # Assegna porte univoche ai reducer
      - "%d:%d"  # Porta univoca per servizio di HeartBeat
    environment:
      - TZ=Europe/Rome
      - REDUCER_PORT=%d
      - HB_PORT=%d
    networks:
      - my-network
`, i, config.ReducerPn+i, config.ReducerPn+i, config.ReducerHbPn+i, config.ReducerHbPn+i, config.ReducerPn+i, config.ReducerHbPn+i)
	}

	outputFile, err := os.Create("docker-compose.yml")
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			log.Fatalf("Something went wrong during container configuration %v", err)
		}
	}(outputFile)

	return tmpl.Execute(outputFile, struct {
		MapperServices  string
		ReducerServices string
	}{
		MapperServices:  mapperServices,
		ReducerServices: reducerServices,
	})
}

func copyConfigurationFile() {
	// Full path of the source file
	sourceRelativePath := "./configuration.json"

	// Full path of the destination directory
	destinationRelativeDir := "./MasterLocale/constants/"

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Calculate the absolute path to the source file
	sourceAbsolutePath := filepath.Join(currentDir, sourceRelativePath)

	// Calculate the absolute path to the destination directory
	destinationAbsoluteDir := filepath.Join(currentDir, destinationRelativeDir)

	// Calculate the absolute path to the destination file
	destinationAbsolutePath := filepath.Join(destinationAbsoluteDir, filepath.Base(sourceAbsolutePath))

	// Open the source file for reading
	sourceFile, err := os.Open(sourceAbsolutePath)
	if err != nil {
		log.Fatalf("Error opening the source file: %v", err)
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			log.Fatalf("Error while closing the source file: %v", err)

		}
	}(sourceFile)

	// Create the destination directory if it doesn't exist
	err = os.MkdirAll(destinationAbsoluteDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating the destination directory: %v", err)
	}

	// Create the destination file
	destinationFile, err := os.Create(destinationAbsolutePath)
	if err != nil {
		log.Fatalf("Error creating the destination file: %v", err)
	}
	defer func(destinationFile *os.File) {
		err := destinationFile.Close()
		if err != nil {
			log.Fatalf("Error while closing the destination file: %v", err)

		}
	}(destinationFile)

	// Copy the content of the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		log.Fatalf("Error copying the file: %v", err)
	}

	log.Printf("File copied successfully to %s", destinationAbsolutePath)
}

func main() {
	configFile, err := os.ReadFile("configuration.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
		return
	}

	err = generateDockerCompose(config)
	if err != nil {
		log.Fatalf("Error generating docker-compose: %v", err)
		return
	}

	if config.Locally == true {
		copyConfigurationFile()
	}

	cmd := exec.Command("docker-compose", "up", "--build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error running docker-compose up: %v", err)
		return
	}

}
