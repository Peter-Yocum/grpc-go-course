package main

import (
	"fmt"
	"log"

	"github.com/Peter-Yocum/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not connect: ", err)
	}
	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)

	fmt.Printf("created client: %f", client)
}
