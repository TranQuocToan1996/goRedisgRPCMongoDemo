package main

import (
	"log"

	"github.com/TranQuocToan1996/redislearn/client"
	"github.com/TranQuocToan1996/redislearn/config"
	"github.com/TranQuocToan1996/redislearn/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "0.0.0.0:8080" // port 8080
)

func main() {
	// insecure.NewCredentials() should not using in PROD cause it turn off TLS
	cfg, _ := config.LoadConfig(".")
	conn, err := grpc.Dial(cfg.GrpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	// Sign Up
	if false {
		signUpUserClient := client.NewSignUpUserClient(conn)
		newUser := &pb.SignUpUserInput{
			Name:            "Jane Smith",
			Email:           "janesmith@gmail.com",
			Password:        "password123",
			PasswordConfirm: "password123",
		}
		signUpUserClient.SignUpUser(newUser)
	}

	// Sign In
	if true {
		signInUserClient := client.NewSignInUserClient(conn)

		credentials := &pb.SignInUserInput{
			Email:    "jamesmith@gmail.com",
			Password: "password123",
		}
		signInUserClient.SignInUser(credentials)
	}

	// Get Me
	if false {

		getMeClient := client.NewGetMeClient(conn)
		id := &pb.GetMeRequest{
			Id: "628cffb91e50302d360c1a2c",
		}
		getMeClient.GetMeUser(id)

	}

}
