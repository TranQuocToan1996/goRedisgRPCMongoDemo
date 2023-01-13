package gapi

import (
	"context"

	"github.com/TranQuocToan1996/redislearn/pb"
	"github.com/TranQuocToan1996/redislearn/services"
	"github.com/TranQuocToan1996/redislearn/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (authServer *AuthServer) SignInUser(ctx context.Context, req *pb.SignInUserInput) (*pb.SignInUserResponse, error) {
	user, err := authServer.userService.FindUserByEmail(req.GetEmail())
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return nil, status.Errorf(codes.InvalidArgument, "Invalid email or password")

		}

		return nil, status.Errorf(codes.Internal, err.Error())

	}

	if !user.Verified {

		return nil, status.Errorf(codes.PermissionDenied, "You are not verified, please verify your email to login")

	}

	if err := utils.Pw.VerifyPassword(user.Password, req.GetPassword()); err != nil {

		return nil, status.Errorf(codes.InvalidArgument, "Invalid email or Password")

	}

	// Generate Tokens
	access_token, err := services.JwtObj.CreateToken(authServer.config.AccessTokenExpiresIn, user.ID.Hex())
	if err != nil {

		return nil, status.Errorf(codes.PermissionDenied, err.Error())

	}

	refresh_token, err := services.JwtObj.CreateRefreshToken(authServer.config.RefreshTokenExpiresIn, user.ID.Hex())
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	res := &pb.SignInUserResponse{
		Status:       "success",
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}

	return res, nil
}
