package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"twitt/rpc"
)

type User struct {
	Username string
	Token    string
	Follows  map[string]*User
	Posts    []*pb.Post
}

var uList map[string]*User

type server struct{}

func (s *server) SignUp(ctx context.Context, in *pb.InfoRequest) (*pb.SuccessReply, error) {
	return &pb.SuccessReply{Success: true}, nil
}

func (s *server) Login(ctx context.Context, in *pb.InfoRequest) (*pb.SuccessReply, error) {
	return &pb.SuccessReply{Success: true}, nil
}

func (s *server) Posting(ctx context.Context, in *pb.PostRequest) (*pb.SuccessReply, error) {
	return &pb.SuccessReply{Success: true}, nil
}

func (s *server) View(ctx context.Context, in *pb.InfoRequest) (*pb.ViewReply, error) {
	return &pb.ViewReply{Success: true}, nil
}

func (s *server) GetList(ctx context.Context, in *pb.InfoRequest) (*pb.ListReply, error) {
	return &pb.ListReply{Success: true}, nil
}

func (s *server) Follow(ctx context.Context, in *pb.FollowingRequest) (*pb.SuccessReply, error) {
	return &pb.SuccessReply{Success: true}, nil
}

func (s *server) UnFollow(ctx context.Context, in *pb.FollowingRequest) (*pb.SuccessReply, error) {
	return &pb.SuccessReply{Success: true}, nil
}


func main() {
	uList = make(map[string]*User)
	lis, err := net.Listen("tcp", ":2233")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTwittServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
