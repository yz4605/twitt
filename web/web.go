package main

import (
	"google.golang.org/grpc"
	"net/http"
	"twitt/pkg/rpc"
	"twitt/pkg/web"
)

func main() {
	conn, _ := grpc.Dial("localhost:2233", grpc.WithInsecure())
	defer conn.Close()
	web.C = pb.NewTwittServiceClient(conn)
	http.HandleFunc("/", web.View)
	http.HandleFunc("/signup", web.SignUp)
	http.HandleFunc("/login", web.Login)
	http.HandleFunc("/post", web.Posting)
	http.HandleFunc("/view", web.View)
	http.HandleFunc("/follow", web.Follow)
	http.HandleFunc("/unfollow", web.UnFollow)
	http.HandleFunc("/logout", web.Logout)
	http.HandleFunc("/config", web.Config)
	http.ListenAndServe(":8080", nil)
}
