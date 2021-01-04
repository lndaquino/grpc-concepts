package main

import (
	"context"
	"fmt"
	"grpc-server/pb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
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

// AddUserVerbose makes a gRPC request to add an user using stream mode
func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "john@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}
		fmt.Println("Status: ", stream.Status, " - ", stream.GetUser().GetName())
	}
}

// AddUsers add users sending a stream of users
func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "w1",
			Name:  "Wesley",
			Email: "wes@wes.com",
		},
		{
			Id:    "w2",
			Name:  "Wesley2",
			Email: "wes2@wes.com",
		},
		{
			Id:    "w3",
			Name:  "Wesley3",
			Email: "wes3@wes.com",
		},
		{
			Id:    "w4",
			Name:  "Wesley4",
			Email: "wes4@wes.com",
		},
		{
			Id:    "w5",
			Name:  "Wesley5",
			Email: "wes5@wes.com",
		},
	}

	// creates stream ready to start streaming
	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// strart streaming users
	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	// stop streaming and receives server response
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

// AddUserStreamBoth send a stream of users and reads the stream of data from server with user status
func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		{
			Id:    "w1",
			Name:  "Wesley",
			Email: "wes@wes.com",
		},
		{
			Id:    "w2",
			Name:  "Wesley2",
			Email: "wes2@wes.com",
		},
		{
			Id:    "w3",
			Name:  "Wesley3",
			Email: "wes3@wes.com",
		},
		{
			Id:    "w4",
			Name:  "Wesley4",
			Email: "wes4@wes.com",
		},
		{
			Id:    "w5",
			Name:  "Wesley5",
			Email: "wes5@wes.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)

			time.Sleep(time.Second * 2)
		}

		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}

		close(wait)
	}()

	<-wait
}
