#!/bin/bash

protoc blog/blogpb/blog.proto --go_out=. --go-grpc_out=.
# python3 -m grpc_tools.protoc -I./blog/blogpb --python_out=./blog/blog_client --grpc_python_out=./blog/blog_client ./blog/blogpb/blog.proto 
