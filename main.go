package main

import (
	"google.golang.org/grpc"
	"net/http"
	"twitt/rpc"
)

var c pb.TwittServiceClient

func main() {
	conn, _ := grpc.Dial("localhost:2233", grpc.WithInsecure())
	defer conn.Close()
	c = pb.NewTwittServiceClient(conn)
	http.HandleFunc("/", View)
	http.HandleFunc("/signup", SignUp)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/post", Posting)
	http.HandleFunc("/view", View)
	http.HandleFunc("/follow", Follow)
	http.HandleFunc("/unfollow", UnFollow)
	http.HandleFunc("/logout", Logout)
	http.ListenAndServe(":8080", nil)
}
