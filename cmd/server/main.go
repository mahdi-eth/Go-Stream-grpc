package main

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/mahdi-eth/go-grpc-streaming/out"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStreamServiceServer
}

func (s server) FetchResponse(in *pb.Request, srv pb.StreamService_FetchResponseServer) error {
	log.Printf("Fetching response for ID: %d", in.Id)

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(i) * time.Second)

		resp := pb.Response{
			Result: fmt.Sprintf("Request #%d for ID:%d", i, in.Id),
		}

		if err := srv.Send(&resp); err != nil {
			log.Printf("Send error: %v", err)
			return err
		}
		log.Printf("Finished request number: %d", i)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterStreamServiceServer(s, &server{})

	log.Println("Starting server on port 50005")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
