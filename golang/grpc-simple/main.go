package main

import (
	"context"
	"log"
	"net"

	pb "alkalax/grpc-simple/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedProviderServer
}

func (s *server) GetStatus(ctx context.Context, in *pb.StatusRequest) (*pb.StatusResponse, error) {
	return &pb.StatusResponse{
		Message: "Alkalax Provider is Online",
		Online:  true,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterProviderServer(srv, &server{})

	log.Printf("server listening at %v", listener.Addr())
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
