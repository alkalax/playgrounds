package main

import (
	pb "alkalax/grpc-simple-stream/proto"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedTaskServiceServer
}

func (s *server) GetTaskUpdates(req *pb.TaskRequest, stream pb.TaskService_GetTaskUpdatesServer) error {
	for i := 1; i <= 5; i++ {
		status := &pb.TaskStatus{
			StatusMessage:      "Processing...",
			ProgressPercentage: int32(i * 20),
		}

		if err := stream.Send(status); err != nil {
			return err
		}

		time.Sleep(time.Second)
	}

	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Printf("failed to create listener: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterTaskServiceServer(srv, &server{})
	reflection.Register(srv)

	log.Printf("server listening at %v", listener.Addr())
	if err = srv.Serve(listener); err != nil {
		log.Printf("failed to server: %v", err)
	}
}
