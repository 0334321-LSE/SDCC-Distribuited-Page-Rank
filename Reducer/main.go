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
	// Obtains reducer port number from environment variable
	reducerPort := os.Getenv("REDUCER_PORT")
	if reducerPort == "" {
		err := errors.New("failed to obtain reducerPort number")
		log.Fatalf("Failed to obtain reducerPort number %v", err)
	}
	reducerListener, err := net.Listen("tcp", ":"+reducerPort)
	if err != nil {
		log.Fatalf("Failed to listen on reducerPort %s %v", reducerPort, err)
	}
	log.Printf("Listening on reducerPort %s \n", reducerPort)

	// Obtains heartbeat port number from environment variable
	heartbeatPort := os.Getenv("HB_PORT")
	if reducerPort == "" {
		err := errors.New("failed to obtain heartbeatPort number")
		log.Fatalf("Failed to obtain heartbeatPort number %v", err)
	}
	hbListener, err := net.Listen("tcp", ":"+heartbeatPort)
	if err != nil {
		log.Fatalf("Failed to listen on heartbeatPort %s %v", heartbeatPort, err)
	}
	log.Printf("Listening on heartbeatPort %s \n", heartbeatPort)

	// Create server for reducer
	reducerServer := reducer.Reducer{}
	// Initialize gRPC server for Reducer
	grpcServer := grpc.NewServer()
	// Register reduce service
	reducer.RegisterReducerServer(grpcServer, &reducerServer)

	// Create server for heartbeat
	hbServer := reducer.ReducerHeartbeat{}
	// Register heartbeat service
	reducer.RegisterReducerHeartbeatServer(grpcServer, &hbServer)

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
		if err := grpcServer.Serve(reducerListener); err != nil {
			log.Fatalf("Failed to serve gRPC on reducerPort %s %v", reducerPort, err)

		}
	}()
	go func() {
		if err := grpcServer.Serve(hbListener); err != nil {
			log.Fatalf("Failed to serve gRPC on heartbeatPort %s %v", heartbeatPort, err)

		}
	}()

	// Keep server running
	log.Println("Server is running.")
	select {}

}
