package main

import (
	"google.golang.org/grpc"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"twitt/pkg/rpc"
	"twitt/pkg/web"
)

var cookie []string

func sendPost(url string, reader io.Reader, fn func(w http.ResponseWriter, r *http.Request), t *testing.T) {
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", cookie[0])
	req.Header.Add("Cookie", cookie[1])
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fn)
	handler.ServeHTTP(rr, req)
	if rr.Header()["Set-Cookie"] != nil {
		cookie = rr.Header()["Set-Cookie"]
	}
}

func sendGet(url string, fn func(w http.ResponseWriter, r *http.Request), t *testing.T) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", cookie[0])
	req.Header.Add("Cookie", cookie[1])
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fn)
	handler.ServeHTTP(rr, req)
	if rr.Header()["Set-Cookie"] != nil {
		cookie = rr.Header()["Set-Cookie"]
	}
	if status := rr.Code; status != http.StatusOK && status != http.StatusFound {
		t.Errorf("Wrong status code: got %v\n",
			status)
	}
}

func init() {
	cookie = make([]string, 2)
	conn, _ := grpc.Dial("localhost:2233", grpc.WithInsecure())
	web.C = pb.NewTwittServiceClient(conn)
}

func TestSignUp(t *testing.T) {
	sendGet("/signup", web.SignUp, t)
	sendPost("/signup", strings.NewReader("username=test1&password=123"), web.SignUp, t)
	sendPost("/signup", strings.NewReader("username=test2&password=123"), web.SignUp, t)
	sendPost("/signup", strings.NewReader("username=test1&password=321"), web.SignUp, t)

}

func TestLogIn(t *testing.T) {
	sendGet("/login", web.Login, t)
	sendPost("/login", strings.NewReader("username=test1&password=123"), web.Login, t)
}

func TestPost(t *testing.T) {
	sendGet("/post", web.Posting, t)
	sendPost("/post", strings.NewReader("message=This is test post"), web.Posting, t)
}

func TestFollow(t *testing.T) {
	sendGet("/follow", web.Follow, t)
	sendPost("/follow", strings.NewReader("test1"), web.Follow, t)
}

func TestView(t *testing.T) {
	sendGet("/view", web.View, t)
}

func TestUnfollow(t *testing.T) {
	sendGet("/unfollow", web.UnFollow, t)
	sendPost("/unfollow", strings.NewReader("test1"), web.UnFollow, t)
}

func TestLogout(t *testing.T) {
	sendGet("/logout", web.Logout, t)
}
