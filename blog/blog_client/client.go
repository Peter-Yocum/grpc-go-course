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

	blog := &blogpb.Blog{
		AuthorId: "Peter",
		Title:    "My First Blog",
		Content:  "Content of first blog",
	}
	response := sendCreateBlogRequest(client, blog)

	//proof trying to read a blog that doesn't exist creates an error but doesn't break program
	junk_blog := sendReadBlogRequest(client, "fake id")
	fmt.Printf("The error for retrieving blog is: %v\n", junk_blog)

	retrieved_blog := sendReadBlogRequest(client, response.Blog.Id)
	fmt.Printf("The retrieved blog is: %v\n", retrieved_blog)

	retrieved_blog.Content = "New content to show update works!"

	//proof updating a blog that doesn't exist creates an error but doesn't break program
	junk_blog = &blogpb.Blog{
		Id:       "62406ccd33ece94df7aac7a8",
		AuthorId: "not a real author",
		Title:    "not a real title",
		Content:  "no real content",
	}
	junk_blog = sendUpdateBlogRequest(client, junk_blog)
	fmt.Printf("The error for updating a blog is: %v\n", junk_blog)

	updated_blog := sendUpdateBlogRequest(client, retrieved_blog)
	fmt.Printf("The updated blog is: %v\n", updated_blog)
}

func sendCreateBlogRequest(client blogpb.BlogServiceClient, blog *blogpb.Blog) *blogpb.CreateBlogResponse {

	fmt.Println("Creating create blog request...")
	in := &blogpb.CreateBlogRequest{
		Blog: blog,
	}
	fmt.Println("Sending create blog request")
	createBlogResponse, err := client.CreateBlog(context.Background(), in)
	if err != nil {
		log.Fatalf("Error when creating blog: %v", err)
	}
	fmt.Printf("Create blog response received: %v\n", createBlogResponse)
	return createBlogResponse
}

func sendReadBlogRequest(client blogpb.BlogServiceClient, id string) *blogpb.Blog {

	res, err := client.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: id,
	})
	if err != nil {
		log.Printf("Error when reading blog: %v", err)
	}
	return res.GetBlog()
}

func sendUpdateBlogRequest(client blogpb.BlogServiceClient, blog *blogpb.Blog) *blogpb.Blog {

	res, err := client.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Printf("Error when updating blog: %v", err)
	}
	return res.GetBlog()
}
