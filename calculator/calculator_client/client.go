package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Peter-Yocum/grpc-go-course/calculator/calculatorpb"
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

	client := calculatorpb.NewCalculatorServiceClient(conn)

	fmt.Printf("created client: %f\n", client)

	sendUnaryRequest(client)

	sendPrimeDecompositionRequest(client)

	SendAverageRequest(client)
}

func sendUnaryRequest(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Unary RPC...")
	req := &calculatorpb.CalculatorRequest{
		Calculation: &calculatorpb.Calculation{
			FirstNumber:  3,
			SecondNumber: 10,
		},
	}
	fmt.Printf("req we will be sending: %v\n", req)
	res, err := client.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling calculate rpc: %v\n", err)
	}
	log.Printf("Response from calculate: %v\n", res.Result)
}

func sendPrimeDecompositionRequest(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do PrimeDecomposition RPC...")
	prime_number := 120
	req := &calculatorpb.PrimeDecompositionRequest{
		PrimeNumber: int64(prime_number),
	}
	resStream, err := client.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while sending prime decomposition rpc: %v\n", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// server done streaming
			break
		}
		if err != nil {
			log.Fatalf("Error while reciving prime decomposition response: %v\n", err)
		}
		log.Printf("One factor of %v is: %v\n", prime_number, msg.GetFactor())
	}
}

func SendAverageRequest(client calculatorpb.CalculatorServiceClient) {

	requests := []*calculatorpb.AverageRequest{
		{
			Number: 5,
		},
		{
			Number: 6,
		},
		{
			Number: 7,
		},
		{
			Number: 8,
		},
	}

	stream, err := client.Average(context.Background())
	if err != nil {
		log.Fatalf("Error when opening up average request stream: %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error when closing stream and receiving response: %v", err)
	}
	fmt.Printf("Received average response: %v\n", response)
}
