package main

import (
	"context"
	"log"
	"net"

	"github.com/hueypark/grpc-tutorial/mmorpg/pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGameServer(grpcServer, &gameServer{})
	log.Println("Server started")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type gameServer struct {
	pb.UnimplementedGameServer
}

func (gameServer) Login(ctx context.Context, loginReq *pb.LoginReq) (*pb.LoginRes, error) {
	log.Printf("Login requested: %v\n", loginReq.Msg)

	return &pb.LoginRes{Msg: loginReq.Msg}, nil
}
