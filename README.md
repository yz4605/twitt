# Final Stage:
```sh
twitt/server
$ go test #do the unit test.
$ go build #compile the server.
$ ./server #keep data server run.

twitt/web
$ go test #do the unit test.
$ go build #compile the web.
$ ./web #start webpage service.

Multi raft servers and Reconfiguration
# Initialize three raft servers
$ ./server --id 1 --cluster http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003 --port 2233
$ ./server --id 2 --cluster http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003 --port 2244
$ ./server --id 3 --cluster http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003 --port 2255

# Add and remove a new raft server
$ curl -XPOST http://localhost:8080/config -d "4=http://127.0.0.1:9004=add"
$ ./server --id 4 --cluster http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003,http://127.0.0.1:9004 --port 2266 --join
$ curl -XPOST http://localhost:8080/config -d "4=http://127.0.0.1:9004=remove"
```
http://localhost:8080/ Access twitt here. The web service is stateless and connected to backend data throught gRPC. Raft data servers are fault tolerant.

URL:
- "/signup" Sign up here and alert if username already exists.
- "/login"  Log in here.
- "/post"  Post twitter here.
- "/view"  View the followings' and own posts.
- "/follow"  Explore all non-folllowed users and make the folllowing.
- "/unfollow"  Un-follow any current following users.
- "/logout"  Log out and clean the cookies.
- "/config" Only accept post request to reconfig the raft clusters.

Structure:  
- twitt/pkg -- All the implementations are here for the reuse purpose.  
- twitt/web -- Web side program.
- twitt/server -- Server side program.
