package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"go.etcd.io/etcd/raft/raftpb"
	"log"
	"sort"
	"twitt/pkg/rpc"
	"context"
)

type User struct {
	Username string
	Token    string
	Follows  map[string]*User
	Posts    []*pb.Post
}

var kvs *kvstore
var conf chan raftpb.ConfChange
var err <-chan error

type Server struct{}

func (s *Server) SignUp(ctx context.Context, in *pb.InfoRequest) (*pb.SuccessReply, error) {
	buf, _ := kvs.Lookup(in.Username)
	if in.Username == "" || buf != "" {
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
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(client); err != nil {
			log.Fatal(err)
		}
		kvs.Propose(username, buf.String())
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
	buf, _ := kvs.Lookup(username)
	if buf == "" {
		return &pb.SuccessReply{Success: false}, nil
	}
	user := &User{}
    if err := gob.NewDecoder(bytes.NewBufferString(buf)).Decode(user); err != nil {
		log.Fatal(err)
	}
	if user != nil && user.Token == token {
		return &pb.SuccessReply{Success: true}, nil
	} else {
		return &pb.SuccessReply{Success: false}, nil
	}
}

func (s *Server) Posting(ctx context.Context, in *pb.PostRequest) (*pb.SuccessReply, error) {
	buf, _ := kvs.Lookup(in.Post.Username)
	if buf == "" {
		return &pb.SuccessReply{Success: false}, nil
	}
	client := &User{}
    if err := gob.NewDecoder(bytes.NewBufferString(buf)).Decode(client); err != nil {
		log.Fatal(err)
	}
	if in.Post.Content == "" || client == nil {
		return &pb.SuccessReply{Success: false}, nil
	}
	client.Posts = append(client.Posts, in.Post)
	var writer bytes.Buffer
	if err := gob.NewEncoder(&writer).Encode(client); err != nil {
		log.Fatal(err)
	}
	kvs.Propose(client.Username, writer.String())
	return &pb.SuccessReply{Success: true}, nil
}

func (s *Server) View(ctx context.Context, in *pb.InfoRequest) (*pb.ViewReply, error) {
	posts := make([]*pb.Post, 0)
	buf, _ := kvs.Lookup(in.Username)
	if buf == "" {
		return &pb.ViewReply{Success: false, Posts: nil}, nil
	}
	client := &User{}
    if err := gob.NewDecoder(bytes.NewBufferString(buf)).Decode(client); err != nil {
		log.Fatal(err)
	}
	if client == nil {
		return &pb.ViewReply{Success: false, Posts: nil}, nil
	}
	for _, user := range client.Follows {
		userbuf, _ := kvs.Lookup(user.Username)
		u := &User{}
	    if err := gob.NewDecoder(bytes.NewBufferString(userbuf)).Decode(u); err != nil {
			log.Fatal(err)
		}
		for _, post := range u.Posts {
			posts = append(posts, post)
		}
	}
	for _, post := range client.Posts {
		posts = append(posts, post)
	}
	sort.Slice(posts, func(i, j int) bool {
	    if posts[i].Time > posts[j].Time {
	        return true
	    }
	    if posts[i].Time < posts[j].Time {
	        return false
	    }
	    return posts[i].Username < posts[j].Username
	})
	return &pb.ViewReply{Success: true, Posts: posts}, nil
}

func (s *Server) GetList(ctx context.Context, in *pb.InfoRequest) (*pb.ListReply, error) {
	buf, _ := kvs.Lookup(in.Username)
	if buf == "" {
		return &pb.ListReply{Success: false, List: nil}, nil
	}
	client := &User{}
    if err := gob.NewDecoder(bytes.NewBufferString(buf)).Decode(client); err != nil {
		log.Fatal(err)
	}
	if client == nil {
		return &pb.ListReply{Success: false, List: nil}, nil
	}
	if in.Instruct == "Follow" {
		list := make([]string, 0)
		for _, buf := range kvs.LookupAll() {
			user := &User{}
		    if err := gob.NewDecoder(bytes.NewBufferString(buf)).Decode(user); err != nil {
				log.Fatal(err)
			}
			// Show all the users who are not followed by client.
			if user.Username != client.Username && client.Follows[user.Username] == nil {
				list = append(list, user.Username)
			}
		}
        sort.Slice(list, func(i, j int) bool {
            return list[i] < list[j]
        })
		return &pb.ListReply{Success: true, List: list}, nil
	} else if in.Instruct == "UnFollow" {
		list := make([]string, 0)
		for _, user := range client.Follows {
			list = append(list, user.Username)
		}
        sort.Slice(list, func(i, j int) bool {
            return list[i] < list[j]
        })
		return &pb.ListReply{Success: true, List: list}, nil
	} else {
		return &pb.ListReply{Success: false, List: nil}, nil
	}
}

func (s *Server) Follow(ctx context.Context, in *pb.FollowingRequest) (*pb.SuccessReply, error) {
	buf, _ := kvs.Lookup(in.Username)
	if buf == "" {
		return &pb.SuccessReply{Success: false}, nil
	}
	client := &User{}
    if err := gob.NewDecoder(bytes.NewBufferString(buf)).Decode(client); err != nil {
		log.Fatal(err)
	}
	if client == nil {
		return &pb.SuccessReply{Success: false}, nil
	}
	buf2, _ := kvs.Lookup(in.Following)
	if buf2 == "" {
		return &pb.SuccessReply{Success: false}, nil
	}
	u := &User{}
    if err := gob.NewDecoder(bytes.NewBufferString(buf2)).Decode(u); err != nil {
		log.Fatal(err)
	}
	if u != nil {
		client.Follows[in.Following] = u
		var writer bytes.Buffer
		if err := gob.NewEncoder(&writer).Encode(client); err != nil {
			log.Fatal(err)
		}
		kvs.Propose(client.Username, writer.String())
		return &pb.SuccessReply{Success: true}, nil
	} else {
		return &pb.SuccessReply{Success: false}, nil
	}
}

func (s *Server) UnFollow(ctx context.Context, in *pb.FollowingRequest) (*pb.SuccessReply, error) {
	buf, _ := kvs.Lookup(in.Username)
	if buf == "" {
		return &pb.SuccessReply{Success: false}, nil
	}
	client := &User{}
    if err := gob.NewDecoder(bytes.NewBufferString(buf)).Decode(client); err != nil {
		log.Fatal(err)
	}
	if client == nil {
		return &pb.SuccessReply{Success: false}, nil
	}
	buf2, _ := kvs.Lookup(in.Following)
	if buf2 == "" {
		return &pb.SuccessReply{Success: false}, nil
	}
	u := &User{}
    if err := gob.NewDecoder(bytes.NewBufferString(buf2)).Decode(u); err != nil {
		log.Fatal(err)
	}
	_, ok := client.Follows[in.Following]
	if u != nil && ok {
    	delete(client.Follows, in.Following)
    	var writer bytes.Buffer
		if err := gob.NewEncoder(&writer).Encode(client); err != nil {
			log.Fatal(err)
		}
		kvs.Propose(client.Username, writer.String())
		return &pb.SuccessReply{Success: true}, nil
	} else {
		return &pb.SuccessReply{Success: false}, nil
	}
}