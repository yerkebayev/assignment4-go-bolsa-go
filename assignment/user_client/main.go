package main

import (
	"context"
	"google.golang.org/grpc"
	"log"

	pb "go-bolsa-go/assignment/user"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	user := &pb.User{
		Name:  "Marat Yerkebayev",
		Email: "210107145@stu.sdu.edu.com",
	}
	addedUser, err := client.AddUser(context.Background(), user)
	if err != nil {
		log.Fatalf("AddUser failed: %v", err)
	}
	log.Printf("User added: %v", addedUser)

	userID := &pb.UserId{Id: addedUser.Id}
	retrievedUser, err := client.GetUser(context.Background(), userID)
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	log.Printf("User retrieved: %v", retrievedUser)

	list, err := client.ListUsers(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("ListUsers failed: %v", err)
	}
	log.Printf("List of users: %v", list.Users)
}
