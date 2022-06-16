package main

import (
	"context"
	"log"

	"github.com/hueypark/grpc-tutorial/mmorpg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:7777", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewGameClient(conn)

	res, err := client.Login(context.Background(), &pb.LoginReq{Msg: "Hello"})
	if err != nil {
		log.Fatalf("fail to login: %v", err)
	}

	log.Printf("Login response: %v\n", res.Msg)
}
