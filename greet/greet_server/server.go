package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Peter-Yocum/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.GreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstname := req.GetGreeting().GetFirstName()
	result := "hello " + firstname
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
