#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=. --go-grpc_out=.
python3 -m grpc_tools.protoc -I./greet/greetpb --python_out=./greet/greet_client --grpc_python_out=./greet/greet_client ./greet/greetpb/greet.proto 
#protoc greet/greetpb/greet.proto --python_out=../greet_client --grpc_python_out=../greet_client