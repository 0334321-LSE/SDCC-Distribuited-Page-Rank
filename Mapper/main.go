package main

import (
	"Mapper/mapper"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	// Obtains mapperPort from environment variable
	mapperPort := os.Getenv("MAPPER_PORT")
	if mapperPort == "" {
		err := errors.New("failed to obtain mapperPort number")
		log.Fatalf("Failed to obtain mapperPort number %v", err)
	}
	mapListener, err := net.Listen("tcp", ":"+mapperPort)
	if err != nil {
		log.Fatalf("Failed to listen on mapperPort %s %v", mapperPort, err)
	}
	log.Printf("Listening on mapperPort %s \n", mapperPort)

	// Obtains heartbeatPort from environment variable
	heartbeatPort := os.Getenv("HB_PORT")
	if heartbeatPort == "" {
		err := errors.New("failed to obtain heartbeatPort number")
		log.Fatalf("Failed to obtain heartbeatPort number %v", err)
	}
	hbListener, err := net.Listen("tcp", ":"+heartbeatPort)
	if err != nil {
		log.Fatalf("Failed to listen on heartbeatPort %s %v", heartbeatPort, err)
	}
	log.Printf("Listening on heartbeatPort %s \n", heartbeatPort)

	// Create server for mapper
	mapperServer := mapper.Mapper{}
	// Initialize gRPC server for Mapper
	grpcServer := grpc.NewServer()
	// Register map service
	mapper.RegisterMapperServer(grpcServer, &mapperServer)

	// Create server for heartbeat
	hbServer := mapper.MapperHeartbeat{}
	// Register heartbeat service
	mapper.RegisterMapperHeartbeatServer(grpcServer, &hbServer)

	// Show exposed services
	func(server *grpc.Server) {
		services := server.GetServiceInfo()
		for keys, service := range services {
			log.Printf("- Nome: %s\n", keys)
			for _, method := range service.Methods {
				log.Printf("  Metodo: %s\n", method.Name)
			}
			log.Println()
		}
	}(grpcServer)

	// Serve both services
	go func() {
		if err := grpcServer.Serve(mapListener); err != nil {
			log.Fatalf("Failed to serve gRPC on mapperPort %s %v", mapperPort, err)

		}
	}()
	go func() {
		if err := grpcServer.Serve(hbListener); err != nil {
			log.Fatalf("Failed to serve gRPC on heartbeatPort %s %v", heartbeatPort, err)
		}
	}()

	// Keep server running
	fmt.Println("Server is running.")
	select {}
}
