package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
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
		t, _ := template.ParseFiles("template/signup.html", "template/headerFooter.html")
		t.Execute(w, nil)
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
	bs, err := ioutil.ReadAll(r.Body)
	username := string(bs)
	if uList[username] != nil || err != nil {
		fmt.Fprint(w, "false")
		return
	} else {
		fmt.Fprint(w, "true")
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/login.html", "template/headerFooter.html")
		t.Execute(w, nil)
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
		t, _ := template.ParseFiles("template/post.html", "template/headerFooter.html")
		t.Execute(w,  nil)
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
	posts := make([]Post, 0)
	for _, user := range client.Follows {
		for _, post := range user.Posts {
			posts = append(posts, *post)
		}
	}
	for _, post := range client.Posts {
		posts = append(posts, *post)
	}
	t, _ := template.ParseFiles("template/view.html", "template/headerFooter.html")
	t.Execute(w, map[string]interface{}{"Posts": posts})
}

func follow(w http.ResponseWriter, r *http.Request) {
	client := validate(w, r, true)
	if client == nil {
		return
	}
	if r.Method == "GET" {
		users := make([]User, 0)
		for _, v := range uList {
			// Show all the users who are not followed by client.
			if v != client && client.Follows[v.UserName] == nil {
				users = append(users, *v)
			}
		}
		t, _ := template.ParseFiles("template/follow.html", "template/headerFooter.html")
		t.Execute(w,  map[string]interface{}{"Users": users})
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
		users := make([]User, 0)
		for _, user := range client.Follows {
			users = append(users, *user)
		}
		t, _ := template.ParseFiles("template/unfollow.html", "template/headerFooter.html")
		t.Execute(w,  map[string]interface{}{"Users": users})
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
