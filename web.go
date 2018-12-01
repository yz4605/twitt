package main

import (
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"twitt/rpc"
)

func validate(w http.ResponseWriter, r *http.Request, flag bool) string {
	username, e := r.Cookie("username")
	isLogin, _ := r.Cookie("isLogin")
	if flag == true && e == nil && isLogin.Value == "true" {
		return username.Value
	}
	expiration := time.Now().Add(-365 * 24 * time.Hour)
	cookie1 := http.Cookie{Name: "username", Value: "", Expires: expiration}
	cookie2 := http.Cookie{Name: "isLogin", Value: "", Expires: expiration}
	http.SetCookie(w, &cookie1)
	http.SetCookie(w, &cookie2)
	http.Redirect(w, r, "/login", 302)
	return ""
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/signup.html", "template/headerFooter.html")
		t.Execute(w, "hidden")
	} else {
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := c.SignUp(ctx, &pb.InfoRequest{Username: username, Password: password})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		if reply.Success {
			http.Redirect(w, r, "/login", 302)
		} else {
			t, _ := template.ParseFiles("template/signup.html", "template/headerFooter.html")
			t.Execute(w, "visible")
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/login.html", "template/headerFooter.html")
		t.Execute(w, "hidden")
	} else {
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := c.Login(ctx, &pb.InfoRequest{Username: username, Password:password})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		if reply.Success {
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie1 := http.Cookie{Name: "username", Value: username, Expires: expiration}
			cookie2 := http.Cookie{Name: "isLogin", Value: "true", Expires: expiration}
			http.SetCookie(w, &cookie1)
			http.SetCookie(w, &cookie2)
			http.Redirect(w, r, "/view", 302)
		} else {
			t, _ := template.ParseFiles("template/login.html", "template/headerFooter.html")
			t.Execute(w, "visible")
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	validate(w, r,false)
}

func Posting(w http.ResponseWriter, r *http.Request) {
	username := validate(w, r,true)
	if username == "" {
		return
	}
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/post.html", "template/headerFooter.html")
		t.Execute(w,  nil)
	} else {
		r.ParseForm()
		message := r.Form["message"][0]
		post := &pb.Post{
			Username: username,
			Content: message,
			Time: time.Now().Format("Mon Jan _2 15:04:05 2006"),
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := c.Posting(ctx, &pb.PostRequest{Post: post})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		if reply.Success {
			http.Redirect(w, r, "/view", 302)
		}
	}
}

func View(w http.ResponseWriter, r *http.Request) {
	username := validate(w, r, true)
	if username == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := c.View(ctx, &pb.InfoRequest{Username:username})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	posts := make([]pb.Post, 0)
	for _, i := range reply.Posts {
		posts = append(posts, *i)
	}

	t, _ := template.ParseFiles("template/view.html", "template/headerFooter.html")
	t.Execute(w, map[string]interface{}{"Posts": posts})
}

func Follow(w http.ResponseWriter, r *http.Request) {
	username := validate(w, r, true)
	if username == "" {
		return
	}
	if r.Method == "GET" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := c.GetList(ctx, &pb.InfoRequest{Username:username, Instruct:"Follow"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		t, _ := template.ParseFiles("template/follow.html", "template/headerFooter.html")
		t.Execute(w, map[string]interface{}{"Usernames": reply.List})
	} else {
		user, _ := ioutil.ReadAll(r.Body)
		target := string(user)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := c.Follow(ctx, &pb.FollowingRequest{Username:username, Following:target})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		if !reply.Success {
			log.Print("Follow Error")
		}
	}
}

func UnFollow(w http.ResponseWriter, r *http.Request) {
	username := validate(w, r, true)
	if username == "" {
		return
	}
	if r.Method == "GET" {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := c.GetList(ctx, &pb.InfoRequest{Username:username, Instruct:"UnFollow"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		t, _ := template.ParseFiles("template/unfollow.html", "template/headerFooter.html")
		t.Execute(w, map[string]interface{}{"Usernames": reply.List})
	} else {
		user, _ := ioutil.ReadAll(r.Body)
		target := string(user)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := c.UnFollow(ctx, &pb.FollowingRequest{Username:username, Following:target})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		if !reply.Success {
			log.Print("UnFollow Error")
		}
	}
}