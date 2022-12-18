package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"github.com/DadaxonOlimov/connect"
	"github.com/DadaxonOlimov/internal"
	"github.com/DadaxonOlimov/proto"
)

func main() {
	log.Println("Starting listening on port 8082")
	port := ":8082"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s", port)
	conn := connect.Connection()
	srv := internal.NewGRPCServer(conn)

	grpcServer := grpc.NewServer()
	proto.RegisterTodoServiceServer(grpcServer, srv)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}