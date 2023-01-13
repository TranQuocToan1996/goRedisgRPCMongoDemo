package main

import (
	"context"
	"log"
	"time"

	"github.com/TranQuocToan1996/redislearn/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "0.0.0.0:8080" // port 8080
)

func main() {
	// insecure.NewCredentials() should not using in PROD cause it turn off TLS
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	newUser := &pb.SignUpUserInput{
		Name:            "Tran Quoc Toan",
		Email:           "tranquoctoan.ce@gmail.com",
		Password:        "password123",
		PasswordConfirm: "password123",
	}

	res, err := client.SignUpUser(ctx, newUser)
	if err != nil {
		log.Fatalf("SignUpUser: %v", err)
	}

	log.Println(res)

}
