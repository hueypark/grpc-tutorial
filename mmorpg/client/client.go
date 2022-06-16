package main

import (
	"context"
	"io"
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

	stream, err := client.Move(context.Background())
	if err != nil {
		log.Fatalf("fail to move: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			movePush, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			} else if err != nil {
				log.Fatalf("fail to recv: %v", err)
			}

			log.Printf("Move pushed: %v\n", movePush.Position)
		}
	}()

	moveReqs := []*pb.MoveReq{{Position: 1}, {Position: 2}, {Position: 3}}
	for _, moveReq := range moveReqs {
		err = stream.Send(moveReq)
		if err != nil {
			log.Fatalf("fail to send: %v", err)
		}
	}
	err = stream.CloseSend()
	if err != nil {
		log.Fatalf("fail to close send: %v", err)
	}

	<-waitc
}
