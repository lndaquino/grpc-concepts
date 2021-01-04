package services

import (
	"context"
	"fmt"
	"grpc-server/pb"
)

// type UserServiceServer interface {
// 	AddUser(context.Context, *User) (*User, error)
// 	mustEmbedUnimplementedUserServiceServer()
// }

// UserService struct models de user service implementation
type UserService struct {
	pb.UnimplementedUserServiceServer
}

// NewUserService returns a new UserService
func NewUserService() *UserService {
	return &UserService{}
}

// AddUser handles grpc request to add user
func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	// insert in database
	fmt.Println(req.Name)

	return &pb.User{
		Id:    "123",
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}
