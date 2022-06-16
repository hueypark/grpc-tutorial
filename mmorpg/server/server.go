package main

import (
	"context"
	"io"
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

func (gameServer) Move(stream pb.Game_MoveServer) error {
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				log.Printf("Move stream closed: %v\n", err)
				close(waitc)
				return
			} else if err != nil {
				close(waitc)
				log.Printf("Move stream err: %v\n", err)
				return
			}

			log.Printf("Move received: %v\n", in.Position)
		}
	}()

	for i := int32(0); i <= 10; i++ {
		err := stream.Send(&pb.MovePush{Position: i})
		if err != nil {
			return err
		}
	}

	<-waitc

	return nil
}
