package gapi

import (
	"github.com/TranQuocToan1996/redislearn/config"
	"github.com/TranQuocToan1996/redislearn/pb"
	"github.com/TranQuocToan1996/redislearn/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	config         config.Config
	userService    services.UserService
	userCollection *mongo.Collection
}

func NewGrpcUserServer(config config.Config, userService services.UserService, userCollection *mongo.Collection) (*UserServer, error) {
	userServer := &UserServer{
		config:         config,
		userService:    userService,
		userCollection: userCollection,
	}

	return userServer, nil
}
