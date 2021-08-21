package authService

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"grpc-auth-mongo/proto"
	"log"
	"time"
)

const (
	usersTable = "users"
)

var (
	ErrIncorrectCredentials = errors.New("incorrect Credentials Supplied")
	ErrUsernameUsed         = errors.New("sorry, the username is already in use")
	ErrEmailUsed         = errors.New("sorry, the email is already in use")
)
type AuthServer struct{
	DB *mongo.Database
	proto.UnimplementedAuthServiceServer
}

func (a AuthServer) Login(_ context.Context, input *proto.LoginRequest) (*proto.AuthResponse, error) {
	login, password := input.GetLogin(), input.GetPassword()
	var user User
	// M is an unordered representation of a BSON documents
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err := a.DB.Collection(usersTable).FindOne(ctx, bson.M{"$or": []bson.M{
		bson.M{"username": login},
		bson.M{"email": login},
	}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("No User Found")
		}
		log.Fatal("MongoDB encountered a fatal err ", err.Error())
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return &proto.AuthResponse{}, ErrIncorrectCredentials
	}
	token := user.UserToJWTToken()
	return &proto.AuthResponse{Token: token}, nil

}

func (a AuthServer) Signup(c context.Context, input *proto.SignupRequest) (*proto.AuthResponse, error) {
	username, email, password := input.GetUsername(), input.GetEmail(), input.GetPassword()
	// ignore validation for simplicity
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	foundUsername, err := a.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: username})
	if err != nil {
		log.Fatal("error when calling username used function")
		return &proto.AuthResponse{}, errors.New("err when trying to check if username used")
	}
	if foundUsername.GetUsed() {
		return &proto.AuthResponse{}, ErrUsernameUsed
	}

	foundEmail, err := a.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: email})
	if err != nil {
		log.Fatal("error when calling email used function")
		return &proto.AuthResponse{}, errors.New("err when trying to check if email used")
	}
	if foundEmail.GetUsed() {
		return &proto.AuthResponse{}, ErrEmailUsed
	}
	
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	_, err = a.DB.Collection(usersTable).InsertOne(ctx, user)
	if err != nil {
		log.Println("encountered an error when inserting a new user", err.Error())
		return &proto.AuthResponse{}, errors.New("error when inserting a  user")
	}
	token := user.UserToJWTToken()
	return &proto.AuthResponse{Token: token },nil

}

func (a AuthServer) UsernameUsed(_ context.Context, input *proto.UsernameUsedRequest) (*proto.UsedResponse, error) {
	found := false
	username := input.GetUsername()
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	var user, emptyUser User
	a.DB.Collection(usersTable).FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if user != emptyUser {
		found = true
	}
	return &proto.UsedResponse{Used: found}, nil
}

func (a AuthServer) EmailUsed(_ context.Context, input *proto.EmailUsedRequest) (*proto.UsedResponse, error) {
	found := false
	email := input.GetEmail()
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	var user, emptyUser User
	a.DB.Collection(usersTable).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if user != emptyUser {
		found = true
	}
	return &proto.UsedResponse{Used: found}, nil
}

func (a AuthServer) AuthUser(_ context.Context, in *proto.AuthUserRequest) (*proto.AuthUserResponse, error) {
	token := in.GetToken()
	user := UserFromJWTToken(token)
	return &proto.AuthUserResponse{ID: user.ID.Hex(), Username: user.Username, Email: user.Email}, nil
}