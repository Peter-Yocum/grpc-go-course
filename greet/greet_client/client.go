package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	//sendUnaryRequest(client)

	//sendStreamingRequest(client)

	sendClientStreamingRequest(client)
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

func sendStreamingRequest(client greetpb.GreetServiceClient) {
	fmt.Println("Starting to do Server Streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Peter",
			LastName:  "Yocum",
		},
	}
	resStream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling greetmanytimes rpc: %v\n", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// reached end of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while collecting results for greetmanytimes rpc: %v\n", err)
		}
		log.Printf("Response from greetmanytimes: %v\n", msg.GetResult())
	}
}

func sendClientStreamingRequest(client greetpb.GreetServiceClient) {

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Peter",
				LastName:  "Yocum",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "reteP",
				LastName:  "mucoY",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Nobody",
				LastName:  "Nemo",
			},
		},
	}

	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while creating streaming client")
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from long greet: %v", err)
	}
	fmt.Printf("Response received for requests: %v\n", response)
}
