package main

import (
	"context"
	"google.golang.org/grpc/codes"
	"net"
	"testing"

	pb "go-bolsa-go/assignment/user"
	"google.golang.org/grpc"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	users []*pb.User
}

func (s *userServiceServer) AddUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	s.users = append(s.users, user)
	user.Id = int32(len(s.users))
	return user, nil
}

func (s *userServiceServer) GetUser(ctx context.Context, userID *pb.UserId) (*pb.User, error) {
	if userID.Id <= 0 && int(userID.Id) > len(s.users) {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}
	return s.users[userID.Id-1], nil
}

func (s *userServiceServer) ListUsers(ctx context.Context, empty *pb.Empty) (*pb.Users, error) {
	return &pb.Users{Users: s.users}, nil
}

func TestServer_AddUser(t *testing.T) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &userServiceServer{})
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Errorf("failed to serve: %v", err)
			return
		}
	}()
	defer grpcServer.Stop()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	user := &pb.User{
		Name:  "Marat Yerkebayev",
		Email: "210107145@stu.sdu.edu.kz",
	}
	addedUser, err := client.AddUser(context.Background(), user)
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}
	if addedUser.Id == 0 {
		t.Error("AddUser did not return a valid user ID")
	}
}

func TestServer_GetUser(t *testing.T) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	srv := &userServiceServer{}
	pb.RegisterUserServiceServer(grpcServer, srv)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Errorf("failed to serve: %v", err)
			return
		}
	}()
	defer grpcServer.Stop()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	user := &pb.User{
		Name:  "Marat Yerkebayev",
		Email: "210107145@stu.sdu.edu.kz",
	}
	addedUser, err := client.AddUser(context.Background(), user)
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	userID := &pb.UserId{Id: addedUser.Id}
	retrievedUser, err := client.GetUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if retrievedUser.Id != addedUser.Id || retrievedUser.Name != addedUser.Name || retrievedUser.Email != addedUser.Email {
		t.Error("Retrieved user does not match added user")
	}
}

func TestServer_ListUsers(t *testing.T) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	srv := &userServiceServer{}
	pb.RegisterUserServiceServer(grpcServer, srv)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Errorf("failed to serve: %v", err)
			return
		}
	}()
	defer grpcServer.Stop()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	users := []*pb.User{
		{Name: "Marat Yerkebayev", Email: "210107145@stu.sdu.edu.com"},
		{Name: "Cristiano Ronaldo", Email: "cr7@gmail.com"},
		{Name: "Toni Kroos", Email: "kr8s@gmail.com"},
	}
	for _, user := range users {
		_, err := client.AddUser(context.Background(), user)
		if err != nil {
			t.Fatalf("AddUser failed: %v", err)
		}
	}

	list, err := client.ListUsers(context.Background(), &pb.Empty{})
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if len(list.Users) != len(users) {
		t.Errorf("Expected %d users, got %d", len(users), len(list.Users))
	}
}
