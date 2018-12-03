# Stage Two:
```sh
twitt/server
$ go test #do the unit test.
$ go build #compile the server.
$ ./server #keep data server run.

twitt/web
$ go test #do the unit test.
$ go build #compile the web.
$ ./web #start webpage service.
```
http://localhost:8080/ Access twitt here and the data is stored in backend so that the web service is stateless. Keep data server running to ensure the communication.  

URL:
- "/signup" Sign up here and alert if username already exists.
- "/login"  Log in here.
- "/post"  Post twitter here.
- "/view"  View the followings' and own posts.
- "/follow"  Explore all non-folllowed users and make the folllowing.
- "/unfollow"  Un-follow any current following users.
- "/logout"  Log out and clean the cookies.

Structure:  
- twitt/pkg -- All the implementations are here for the reuse purpose.  
- twitt/web -- Web side main function.  
- twitt/server -- Server side main function.  
