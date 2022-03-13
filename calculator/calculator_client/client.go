package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"time"

	"github.com/Peter-Yocum/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

	//sendUnaryRequest(client)

	//sendPrimeDecompositionRequest(client)

	//sendAverageRequest(client)

	//sendFindMaximumRequest(client)

	sendSquareRootRequest(client)
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

func sendAverageRequest(client calculatorpb.CalculatorServiceClient) {

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

func sendFindMaximumRequest(client calculatorpb.CalculatorServiceClient) {

	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error when opening up find maximum stream: %v", err)
	}

	numbers_to_send := []int{1, -1, 3, -4, 10, 5, 101, 52, 1010}
	waitchannel := make(chan struct{})
	maximum := math.Inf(-1)

	go func() {
		for _, number := range numbers_to_send {
			fmt.Printf("Sending find maximum request for num: %v\n", number)
			time.Sleep(100 * time.Millisecond)
			req_to_send := &calculatorpb.FindMaximumRequest{
				NextNumber: float32(number),
			}
			err := stream.Send(req_to_send)
			if err != nil {
				log.Fatalf("Error while trying to stream find maximum request: %v\n", err)
			}
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(waitchannel)
				break
			}
			if err != nil {
				log.Fatalf("Error while trying to stream find maximum response: %v\n", err)
			}
			if response.GetCurrentMax() > float32(maximum) {
				maximum = float64(response.GetCurrentMax())
				fmt.Printf("Foud new maximum: %v\n", maximum)
			}
		}
	}()

	<-waitchannel
}

func sendSquareRootRequest(client calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do square root RPC...")
	req := &calculatorpb.SquareRootRequest{
		Number: 1000,
	}
	fmt.Printf("number we will be sending: %v\n", req.GetNumber())
	res, err := client.SquareRoot(context.Background(), req)
	if err != nil {
		processSquareRootError(err)
	}
	log.Printf("Response from square root of num: %v is: %v\n", req.GetNumber(), res.GetNumberRoot())

	req = &calculatorpb.SquareRootRequest{
		Number: -1000,
	}
	fmt.Printf("number we will be sending: %v\n", req.GetNumber())
	res, err = client.SquareRoot(context.Background(), req)
	if err != nil {
		processSquareRootError(err)
	}
	log.Printf("Response from square root of num: %v is: %v\n", req.GetNumber(), res.GetNumberRoot())
}

func processSquareRootError(err error) {
	respErr, ok := status.FromError(err)
	if ok {
		// my error from grpc
		fmt.Println(respErr.Message())
		fmt.Println(respErr.Code())
		if respErr.Code() == codes.InvalidArgument {
			fmt.Println("We sent a negative number!")
		}
	} else {
		// framework error/problem
		log.Fatalf("Error while calling calculate rpc: %v\n", err)
	}
}
