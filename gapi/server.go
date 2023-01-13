package gapi

import (
	"context"
	"html/template"
	"strings"
	"time"

	"github.com/TranQuocToan1996/redislearn/config"
	"github.com/TranQuocToan1996/redislearn/models"
	"github.com/TranQuocToan1996/redislearn/pb"
	"github.com/TranQuocToan1996/redislearn/services"
	"github.com/TranQuocToan1996/redislearn/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAuthServiceServer

	config         config.Config
	authService    services.AuthService
	userService    services.UserService
	userCollection *mongo.Collection
	temp           *template.Template
}

func NewGrpcServer(config config.Config, authService services.AuthService,
	userService services.UserService, userCollection *mongo.Collection, temp *template.Template) (*Server, error) {

	server := &Server{
		config:         config,
		authService:    authService,
		userService:    userService,
		userCollection: userCollection,
		temp:           temp,
	}

	return server, nil
}

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.GenericResponse, error) {
	code := req.GetVerificationCode()

	verificationCode := utils.Encode(code)

	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}, {Key: "updated_at", Value: time.Now()}}}, {Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}
	result, err := server.userCollection.UpdateOne(ctx, query, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.PermissionDenied, "Could not verify email address")
	}

	res := &pb.GenericResponse{
		Status:  "success",
		Message: "Email verified successfully",
	}
	return res, nil
}

func (server *Server) SignUpUser(ctx context.Context, req *pb.SignUpUserInput) (*pb.GenericResponse, error) {
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}

	user := models.SignUpInput{
		Name:            req.GetName(),
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}

	newUser, err := server.authService.SignUpUser(&user)

	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())

		}
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	// Generate Verification Code
	code := utils.RandStringRunes(20)

	verificationCode := utils.Encode(code)

	// Update User in Database
	server.userService.UpdateUserById(newUser.ID.Hex(), "verificationCode", verificationCode)

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[0]
	}

	emailData := utils.EmailData{
		URL:       server.config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(newUser, &emailData, server.temp, "verificationCode.html")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "There was an error sending email: %s", err.Error())

	}

	message := "We sent an email with a verification code to " + newUser.Email

	res := &pb.GenericResponse{
		Status:  "success",
		Message: message,
	}
	return res, nil
}
