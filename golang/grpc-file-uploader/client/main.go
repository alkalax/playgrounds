package main

import (
	pb "alkalax/grpc-file-uploader/proto"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.UploadFile(ctx)
	if err != nil {
		log.Fatalf("failed to initiate stream: %v", err)
	}

	data := []byte("Very large file here, we're simulating file upload in smaller chunks...")
	chunkSize := 2
	for i := 0; i < len(data); i += chunkSize {
		end := min(i+chunkSize, len(data))

		if err := stream.Send(&pb.UploadRequest{
			DataChunk: data[i:end],
		}); err != nil {
			log.Fatalf("failed to send chunk [%d,%d): %v", i, end, err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to close stream: %v", err)
	}
	log.Printf("Upload finished - ID: %s, Total: %d bytes", res.FileId, res.TotalBytesReceived)
}
