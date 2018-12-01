package main

import (
	"crypto/sha256"
	"context"
	"encoding/base64"
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
	if in.Username == "" || uList[in.Username] != nil {
		// This username is not permitted.
		return &pb.SuccessReply{Success: false}, nil
	} else {
		username := in.Username
		password := in.Password
		h := sha256.New()
		h.Write([]byte(password))
		client := new(User)
		client.Username = username
		client.Follows = make(map[string]*User)
		client.Token = base64.StdEncoding.EncodeToString(h.Sum(nil))
		uList[username] = client
		return &pb.SuccessReply{Success: true}, nil
	}
}

func (s *server) Login(ctx context.Context, in *pb.InfoRequest) (*pb.SuccessReply, error) {
	username := in.Username
	password := in.Password
	h := sha256.New()
	h.Write([]byte(password))
	token := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// Check if the username exists and password is correct.
	if uList[username] != nil && uList[username].Token == token {
		return &pb.SuccessReply{Success: true}, nil
	} else {
		return &pb.SuccessReply{Success: false}, nil
	}
}

func (s *server) Posting(ctx context.Context, in *pb.PostRequest) (*pb.SuccessReply, error) {
	client := uList[in.Post.Username]
	client.Posts = append(client.Posts, in.Post)
	return &pb.SuccessReply{Success: true}, nil
}

func (s *server) View(ctx context.Context, in *pb.InfoRequest) (*pb.ViewReply, error) {
	posts := make([]*pb.Post, 0)
	client := uList[in.Username]
	if client == nil {
		return &pb.ViewReply{Success: false, Posts: nil}, nil
	}
	for _, user := range client.Follows {
		for _, post := range user.Posts {
			posts = append(posts, post)
		}
	}
	for _, post := range client.Posts {
		posts = append(posts, post)
	}
	return &pb.ViewReply{Success: true, Posts: posts}, nil
}

func (s *server) GetList(ctx context.Context, in *pb.InfoRequest) (*pb.ListReply, error) {
	client := uList[in.Username]
	if in.Instruct == "Follow" {
		list := make([]string, 0)
		for _, user := range uList {
			// Show all the users who are not followed by client.
			if user != client && client.Follows[user.Username] == nil {
				list = append(list, user.Username)
			}
		}
		return &pb.ListReply{Success: true, List: list}, nil
	} else {
		list := make([]string, 0)
		for _, user := range client.Follows {
			list = append(list, user.Username)
		}
		return &pb.ListReply{Success: true, List: list}, nil
	}
}

func (s *server) Follow(ctx context.Context, in *pb.FollowingRequest) (*pb.SuccessReply, error) {
	client := uList[in.Username]
	u := uList[in.Following]
	if u != nil {
		client.Follows[in.Following] = u
	}
	return &pb.SuccessReply{Success: true}, nil
}

func (s *server) UnFollow(ctx context.Context, in *pb.FollowingRequest) (*pb.SuccessReply, error) {
	client := uList[in.Username]
	u := uList[in.Following]
	if u != nil {
		delete(client.Follows, in.Following)
	}
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
