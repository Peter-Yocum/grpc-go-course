package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Peter-Yocum/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	opts := grpc.WithInsecure()

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v\n", err)
	}
	defer conn.Close()

	client := blogpb.NewBlogServiceClient(conn)

	fmt.Printf("created client: %f\n", client)

	fmt.Println("Creating blog object and request...")
	blog := &blogpb.Blog{
		AuthorId: "Peter",
		Title:    "My First Blog",
		Content:  "Content of first blog",
	}
	in := &blogpb.CreateBlogRequest{
		Blog: blog,
	}
	fmt.Println("Sending create blog request")
	createBlogResponse, err := client.CreateBlog(context.Background(), in)
	if err != nil {
		log.Fatalf("Error when creating blog: %v", err)
	}
	fmt.Printf("Create blog response received: %v", createBlogResponse)
}

// func sendCreateBlogRequest(client greetpb.GreetServiceClient) {
// 	fmt.Println("Starting to do Unary RPC...")
// 	req := &greetpb.GreetRequest{
// 		Greeting: &greetpb.Greeting{
// 			FirstName: "Peter",
// 			LastName:  "Yocum",
// 		},
// 	}
// 	res, err := client.Greet(context.Background(), req)
// 	if err != nil {
// 		log.Fatalf("Error while calling greet rpc: %v\n", err)
// 	}
// 	log.Printf("Response from greet: %v\n", res.Result)
// }
