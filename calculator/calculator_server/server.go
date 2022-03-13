package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/Peter-Yocum/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.CalculatorServiceServer
}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Calculate function was invoked with %v\n", req)
	first_number := req.GetCalculation().GetFirstNumber()
	second_number := req.GetCalculation().GetSecondNumber()
	res := &calculatorpb.CalculatorResponse{
		Result: first_number + second_number,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Prime Decomposition function was invoked with %v\n", req)
	prime := req.GetPrimeNumber()
	factor := int64(2)
	for prime > 1 {
		if prime%factor == 0 {
			res := &calculatorpb.PrimeDecompositionResponse{
				Factor: factor,
			}
			stream.Send(res)
			prime = prime / factor
		} else {
			factor++
		}

	}
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	fmt.Print("Starting average calculation\n")
	total := float32(0)
	num_req := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			result := total / float32(num_req)
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Result: float32(result),
			})
		}
		if err != nil {
			log.Fatalf("error when trying to receive stream message in average: %v", err)
		}
		total += float32(req.GetNumber())
		num_req++
		fmt.Printf("Received another number to average, running total: %v, total numbers received: %v\n", total, num_req)
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Print("Starting Find Maximum\n")
	current_max := float32(math.Inf(-1))
	for {
		req, recv_err := stream.Recv()
		if recv_err == io.EOF {
			return nil
		}
		if recv_err != nil {
			log.Fatalf("error when trying to receive stream message in find maximum: %v", recv_err)
			return recv_err
		}
		new_number := req.GetNextNumber()
		if new_number > current_max {
			current_max = new_number
		}
		send_err := stream.SendMsg(&calculatorpb.FindMaximumResponse{
			CurrentMax: current_max,
		})
		if send_err != nil {
			log.Fatalf("error when trying to send stream message in find maximum: %v", send_err)
			return send_err
		}
	}
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
