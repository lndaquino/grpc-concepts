package main

import (
	"context"
	"fmt"
	"grpc-server/pb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	AddUser(client)
}

// AddUser makes a gRPC request to add an user
func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "john@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}
