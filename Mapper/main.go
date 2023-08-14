package main

import (
	"Mapper/mapper"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	// Obtains assigned port on Docker Compose from environment variable
	port := os.Getenv("EXPOSED_PORT")
	if port == "" {
		err := errors.New("failed to obtain port number")
		log.Fatalf("Failed to obtain port number %v", err)
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s %v", port, err)
	}
	log.Printf("Listening on port %s \n", port)

	mapperServer := mapper.Mapper{}

	// Initialize gRPC server for Mapper
	grpcServer := grpc.NewServer()

	mapper.RegisterMapperServer(grpcServer, &mapperServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC on port %s %v", port, err)

	}
}
