package main

import (
	"context"
	"fmt"
	"log"
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
