package main

import (
	"Reducer/reducer"
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

	reducerServer := reducer.Reducer{}

	// Initialize gRPC server for Reducer
	grpcServer := grpc.NewServer()

	reducer.RegisterReducerServer(grpcServer, &reducerServer)

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

	// Serve
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC on port %s %v", port, err)

	}

}
