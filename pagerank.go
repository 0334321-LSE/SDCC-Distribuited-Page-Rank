package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
)

type Config struct {
	NumMapper          int     `json:"num_mapper"`
	NumReducer         int     `json:"num_reducer"`
	MapperPN           int     `json:"mapper_pn"`
	MapperHeartBeatPN  int     `json:"mapper_hb_pn"`
	ReducerPN          int     `json:"reducer_pn"`
	ReducerHeartBeatPN int     `json:"reducer_hb_pn"`
	DampingFactor      float64 `json:"damping_factor"`
	Epsilon            float64 `json:"epsilon"`
	GraphPath          string  `json:"graph_path"`
	MaxIteration       int     `json:"max_iteration"`
	NumNodes           int     `json:"num_nodes"`
	EdgesToAttach      int     `json:"edges_to_attach"`
	Seed               int64   `json:"seed"`
}

func generateDockerCompose(config Config) error {
	composeTemplate := `
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
`, i, config.MapperPN+i, config.MapperPN+i, config.MapperHeartBeatPN+i, config.MapperHeartBeatPN+i, config.MapperPN+i, config.MapperHeartBeatPN+i)
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
`, i, config.ReducerPN+i, config.ReducerPN+i, config.ReducerHeartBeatPN+i, config.ReducerHeartBeatPN+i, config.ReducerPN+i, config.ReducerHeartBeatPN+i)
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

func main() {
	configFile, err := os.ReadFile("configuration.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println("Error parsing config:", err)
		return
	}

	err = generateDockerCompose(config)
	if err != nil {
		fmt.Println("Error generating docker-compose:", err)
		return
	}

	cmd := exec.Command("docker-compose", "up", "--build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running docker-compose up:", err)
		return
	}
}
