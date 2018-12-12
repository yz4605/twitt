package server

import (
	"crypto/sha256"
	"encoding/base64"
	"go.etcd.io/etcd/raft/raftpb"
	"twitt/pkg/rpc"
	"context"
)

type User struct {
	Username string
	Token    string
	Follows  map[string]*User
	Posts    []*pb.Post
}

var UList map[string]*User
var kvs *kvstore
var conf chan raftpb.ConfChange
var err <-chan error

type Server struct{}

func init() {
	UList = make(map[string]*User)
}

func (s *Server) SignUp(ctx context.Context, in *pb.InfoRequest) (*pb.SuccessReply, error) {
	if in.Username == "" || UList[in.Username] != nil {
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
		UList[username] = client
		return &pb.SuccessReply{Success: true}, nil
	}
}

func (s *Server) Login(ctx context.Context, in *pb.InfoRequest) (*pb.SuccessReply, error) {
	username := in.Username
	password := in.Password
	h := sha256.New()
	h.Write([]byte(password))
	token := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// Check if the username exists and password is correct.
	if UList[username] != nil && UList[username].Token == token {
		return &pb.SuccessReply{Success: true}, nil
	} else {
		return &pb.SuccessReply{Success: false}, nil
	}
}

func (s *Server) Posting(ctx context.Context, in *pb.PostRequest) (*pb.SuccessReply, error) {
	client := UList[in.Post.Username]
	if in.Post.Content == "" || client == nil {
		return &pb.SuccessReply{Success: false}, nil
	}
	client.Posts = append(client.Posts, in.Post)
	return &pb.SuccessReply{Success: true}, nil
}

func (s *Server) View(ctx context.Context, in *pb.InfoRequest) (*pb.ViewReply, error) {
	posts := make([]*pb.Post, 0)
	client := UList[in.Username]
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

func (s *Server) GetList(ctx context.Context, in *pb.InfoRequest) (*pb.ListReply, error) {
	client := UList[in.Username]
	if client == nil {
		return &pb.ListReply{Success: false, List: nil}, nil
	}
	if in.Instruct == "Follow" {
		list := make([]string, 0)
		for _, user := range UList {
			// Show all the users who are not followed by client.
			if user != client && client.Follows[user.Username] == nil {
				list = append(list, user.Username)
			}
		}
		return &pb.ListReply{Success: true, List: list}, nil
	} else if in.Instruct == "UnFollow" {
		list := make([]string, 0)
		for _, user := range client.Follows {
			list = append(list, user.Username)
		}
		return &pb.ListReply{Success: true, List: list}, nil
	} else {
		return &pb.ListReply{Success: false, List: nil}, nil
	}
}

func (s *Server) Follow(ctx context.Context, in *pb.FollowingRequest) (*pb.SuccessReply, error) {
	client := UList[in.Username]
	if client == nil {
		return &pb.SuccessReply{Success: false}, nil
	}
	u := UList[in.Following]
	if u != nil {
		client.Follows[in.Following] = u
		return &pb.SuccessReply{Success: true}, nil
	} else {
		return &pb.SuccessReply{Success: false}, nil
	}
}

func (s *Server) UnFollow(ctx context.Context, in *pb.FollowingRequest) (*pb.SuccessReply, error) {
	client := UList[in.Username]
	if client == nil {
		return &pb.SuccessReply{Success: false}, nil
	}
	u := UList[in.Following]
	_, ok := client.Follows[in.Following]
	if u != nil && ok {
    	delete(client.Follows, in.Following)
		return &pb.SuccessReply{Success: true}, nil
	} else {
		return &pb.SuccessReply{Success: false}, nil
	}
}