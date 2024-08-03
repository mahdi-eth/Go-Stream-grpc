package main

import (
	"context"
	"io"
	"log"

	pb "github.com/mahdi-eth/go-grpc-streaming/out"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(":50005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)

	in := &pb.Request{Id: 1}

	stream, err := client.FetchResponse(context.Background(), in)
	if err != nil {
		log.Fatalf("open stream error: %v", err)
	}

	done := make(chan bool)

	go func() {
		defer close(done)
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("cannot receive message: %v", err)
			}
			log.Printf("Response received: %s", resp.Result)
		}
	}()

	<-done
	log.Println("Finished receiving messages")
}
