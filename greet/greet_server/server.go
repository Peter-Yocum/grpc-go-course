package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Peter-Yocum/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}
