#!/bin/bash

protoc calculator/calculatorpb/calculator.proto --go_out=. --go-grpc_out=.