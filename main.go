package main

import (
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	UserName string
	Token    string
	Follows  map[string]*User
	Posts    []*Post
}

type Post struct {
	Msg      string
	Time     string
	UserName string
}

var uList map[string]*User

func validate(w http.ResponseWriter, r *http.Request, flag bool) *User {
	cookie1, e := r.Cookie("username")
	cookie2, _ := r.Cookie("token")
	if flag == true && e == nil && uList[cookie1.Value] != nil && uList[cookie1.Value].Token == cookie2.Value {
		return uList[cookie1.Value]
	}
	expiration := time.Now().Add(-365 * 24 * time.Hour)
	cookie3 := http.Cookie{Name: "username", Value: "", Expires: expiration}
	cookie4 := http.Cookie{Name: "token", Value: "", Expires: expiration}
	http.SetCookie(w, &cookie3)
	http.SetCookie(w, &cookie4)
	http.Redirect(w, r, "/login", 302)
	return nil
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
	} else {
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		// User is registered
		if uList[username] != nil {
			return
		}
		h := sha256.New()
		h.Write([]byte(password))
		client := new(User)
		client.Follows = make(map[string]*User)
		client.UserName = username
		client.Token = base64.StdEncoding.EncodeToString(h.Sum(nil))
		uList[username] = client
		http.Redirect(w, r, "/login", 302)
	}
}

func checkUserName(w http.ResponseWriter, r *http.Request) {
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
	} else {
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		h := sha256.New()
		h.Write([]byte(password))
		token := base64.StdEncoding.EncodeToString(h.Sum(nil))
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie1 := http.Cookie{Name: "username", Value: username, Expires: expiration}
		cookie2 := http.Cookie{Name: "token", Value: token, Expires: expiration}
		http.SetCookie(w, &cookie1)
		http.SetCookie(w, &cookie2)
		http.Redirect(w, r, "/view", 302)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	validate(w, r,false)
}

func post(w http.ResponseWriter, r *http.Request) {
	client := validate(w, r,true)
	if client == nil {
		return
	}
	if r.Method == "GET" {
	} else {
		r.ParseForm()
		message := r.Form["message"][0]
		post := Post{
			Msg: message,
			Time: time.Now().Format("Mon Jan _2 15:04:05 2006"),
			UserName: client.UserName,
		}
		client.Posts = append(client.Posts, &post)
		http.Redirect(w, r, "/view", 302)
	}
}

func view(w http.ResponseWriter, r *http.Request) {
	client := validate(w, r, true)
	if client == nil {
		return
	}
}

func follow(w http.ResponseWriter, r *http.Request) {
	client := validate(w, r, true)
	if client == nil {
		return
	}
	if r.Method == "GET" {
	} else {
		user, _ := ioutil.ReadAll(r.Body)
		username := string(user)
		u := uList[username]
		if u != nil {
			client.Follows[username] = u
		}
	}
}

func unfollow(w http.ResponseWriter, r *http.Request) {
	client := validate(w, r, true)
	if client == nil {
		return
	}
	if r.Method == "GET" {
	} else {
		user, _ := ioutil.ReadAll(r.Body)
		username := string(user)
		u := uList[username]
		if u != nil {
			delete(client.Follows, username)
		}
	}
}

func main() {
	uList = make(map[string]*User)
	http.HandleFunc("/", view)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/checkUserName", checkUserName)
	http.HandleFunc("/login", login)
	http.HandleFunc("/post", post)
	http.HandleFunc("/view", view)
	http.HandleFunc("/follow", follow)
	http.HandleFunc("/unfollow", unfollow)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe(":8080", nil)
}
