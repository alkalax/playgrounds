package main

import (
	pb "alkalax/grpc-file-uploader/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

func (s *server) UploadFiles(stream pb.FileService_UploadFileServer) error {
	var totalSize int32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Stream finished, received %d bytes", totalSize)

			return stream.SendAndClose(&pb.UploadResponse{
				FileId:             "unique-uuid-123",
				TotalBytesReceived: totalSize,
			})
		}

		if err != nil {
			return err
		}

		chunkSize := int32(len(req.DataChunk))
		totalSize += chunkSize
		log.Printf("Received chunk of size: %d bytes", chunkSize)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("error while listening: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterFileServiceServer(srv, &server{})
	reflection.Register(srv)

	log.Printf("server listening at %v", listener.Addr())
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
