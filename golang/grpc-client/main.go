package main

import (
	"context"
	"log"
	"time"

	pb "alkalax/grpc-client/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewProviderClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.Add(ctx, &pb.Nums{A: 2, B: 3})
	if err != nil {
		log.Fatalf("calling add: %v", err)
	}

	log.Printf("result: %v", res)
}
