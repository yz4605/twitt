package main

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var c []string

func sendPost(url string, reader io.Reader, fn func(w http.ResponseWriter, r *http.Request), t *testing.T) {
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", c[0])
	req.Header.Add("Cookie", c[1])
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fn)
	handler.ServeHTTP(rr, req)
	if rr.Header()["Set-Cookie"] != nil {
		c = rr.Header()["Set-Cookie"]
	}
}

func TestSignUp(t *testing.T) {
	uList = make(map[string]*User)
	c = make([]string, 2)

	sendPost("/signup", strings.NewReader("username=test1&password=123"), signUp, t)
	sendPost("/signup", strings.NewReader("username=test2&password=123"), signUp, t)
	sendPost("/signup", strings.NewReader("username=test1&password=321"), signUp, t)

	h := sha256.New()
	h.Write([]byte("123"))
	token := base64.StdEncoding.EncodeToString(h.Sum(nil))
	eq := uList["test1"].UserName == "test1" && uList["test1"].Token == token
	if len(uList) != 2 || !eq {
		t.Errorf("Sign up fail")
	}
}

func TestLogIn(t *testing.T) {
	h := sha256.New()
	h.Write([]byte("123"))
	token := base64.StdEncoding.EncodeToString(h.Sum(nil))
	sendPost("/login", strings.NewReader("username=test1&password=123"), login, t)
	if uList["test1"].Token != token {
		t.Errorf("Login fail")
	}
}

func TestPost(t *testing.T) {
	sendPost("/post", strings.NewReader("message=This is test post"), post, t)
	if len(uList["test1"].Posts) != 1 {
		t.Errorf("Post fail")
	}
}

func TestLogIn2(t *testing.T) {
	h := sha256.New()
	h.Write([]byte("123"))
	token := base64.StdEncoding.EncodeToString(h.Sum(nil))
	sendPost("/login", strings.NewReader("username=test2&password=123"), login, t)
	if uList["test2"].Token != token {
		t.Errorf("Login fail")
	}
}

func TestFollow(t *testing.T) {
	sendPost("/follow", strings.NewReader("test1"), follow, t)
	if len(uList["test2"].Follows) != 1 {
		t.Errorf("Follow fail")
	}
}

func TestVew(t *testing.T) {
	req, err := http.NewRequest("GET", "/view", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", c[0])
	req.Header.Add("Cookie", c[1])
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(view)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUnfollow(t *testing.T) {
	sendPost("/unfollow", strings.NewReader("test1"), unfollow, t)
	if len(uList["test2"].Follows) != 0 {
		t.Errorf("unFollow fail")
	}
}

func TestLogout(t *testing.T) {
	req, err := http.NewRequest("GET", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logout)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Wrong status code: got %v want %v",
			status, http.StatusFound)
	}
}
