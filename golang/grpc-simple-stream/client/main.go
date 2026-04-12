package main

import (
	pb "alkalax/grpc-simple-stream/proto"
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.GetTaskUpdates(ctx, &pb.TaskRequest{TaskId: "123"})
	if err != nil {
		log.Fatalf("failed to initiate stream: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Update: %s (%d%%)", msg.StatusMessage, msg.GetProgressPercentage())
	}

	log.Print("End of stream")
}
