package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-auth-mongo/proto"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}

	// anon struct
	user := struct{
		Username string
		Email string
		Password string
	}{
		Username: "William",
		Email:    "will@example.com",
		Password: "Kindlyletmein",
	}

	client := proto.NewAuthServiceClient(conn)
	resp, err := client.Signup(context.Background(), &proto.SignupRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		log.Println("encountered an error when signing-up")
	}

	token := resp.GetToken()
	fmt.Printf("Token: %s\n\n", token)

	resp, err = client.Login(context.Background(), &proto.LoginRequest{
		Login:    "will@example.com",
		Password: "Kindlyletmein",
	})
	if err != nil {
		log.Println("encountered an error when logging-in")
	}
	token = resp.GetToken()
	fmt.Printf("Token: %s\n\n", token)
}