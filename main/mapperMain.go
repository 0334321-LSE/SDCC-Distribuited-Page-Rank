package main

import (
	"PageRank/mapper"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000 %v", err)
	}
	log.Printf("Listening on port 9000 \n")
	mapperServer := mapper.Mapper{}
	// Initialize gRPC server of Mapper
	grpcServer := grpc.NewServer()

	mapper.RegisterMapperServer(grpcServer, &mapperServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC on port 9000 %v", err)

	}
}
