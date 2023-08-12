package main

import (
	"Reducer/reducer"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("Failed to listen on port 9001 %v", err)
	}
	log.Printf("Listening on port 9001 \n")

	reducerServer := reducer.Reducer{}

	// Initialize gRPC server for Reducer
	grpcServer := grpc.NewServer()

	reducer.RegisterReducerServer(grpcServer, &reducerServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC on port 9001 %v", err)

	}
}
