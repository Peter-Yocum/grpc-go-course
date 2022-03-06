package main

import (
	"context"
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
		log.Fatalf("could not connect: %v\n", err)
	}
	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)

	fmt.Printf("created client: %f\n", client)

	sendUnaryRequest(client)
}

func sendUnaryRequest(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Peter",
			LastName:  "Yocum",
		},
	}
	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling greet rpc: %v\n", err)
	}
	log.Printf("Response from greet: %v\n", res.Result)
}
