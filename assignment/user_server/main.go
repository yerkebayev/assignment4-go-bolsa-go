package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"net"

	pb "go-bolsa-go/assignment/user"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	users []*pb.User
}

func (s *userServiceServer) AddUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	userID := int32(len(s.users) + 1)
	user.Id = userID
	s.users = append(s.users, user)
	return user, nil
}

func (s *userServiceServer) GetUser(ctx context.Context, userID *pb.UserId) (*pb.User, error) {
	for _, u := range s.users {
		if u.Id == userID.Id {
			return u, nil
		}
	}
	return nil, grpc.Errorf(codes.NotFound, "user with ID %v not found", userID.Id)
}

func (s *userServiceServer) ListUsers(ctx context.Context, empty *pb.Empty) (*pb.Users, error) {
	return &pb.Users{Users: s.users}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &userServiceServer{})
	log.Println("gRPC server is running on port 50051...")
	err1 := grpcServer.Serve(lis)
	if err1 != nil {
		log.Fatalf("failed to serve: %v", err1)
	}
}
