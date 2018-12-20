# Final Stage:
```sh
twitt/server
$ go test #temp file will be removed.
$ go build #compile the server.

Stand-alone mode
$ ./server #start the data server with default argument as "--id 1 --cluster http://127.0.0.1:9001 --port 2233".

#Shutdown the running server and clean twitt/server/storage to switch the mode, otherwise follow reconfiguration instruction.

Multi-raft mode
$ ./server --id 1 --cluster http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003 --port 2233
$ ./server --id 2 --cluster http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003 --port 2244
$ ./server --id 3 --cluster http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003 --port 2255

Reconfiguration
$ curl -XPOST http://localhost:8080/config -d "4=http://127.0.0.1:9004=add"
$ ./server --id 4 --cluster http://127.0.0.1:9001,http://127.0.0.1:9002,http://127.0.0.1:9003,http://127.0.0.1:9004 --port 2266 --join
$ curl -XPOST http://localhost:8080/config -d "4=http://127.0.0.1:9004=remove"

twitt/web
$ go test #temp data will remain.
$ go build #compile the web.
$ ./web #start webpage service with default argument as "--port 2233".
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
- "/config" Only accept post request to reconfig the raft clusters and get request is not allowed.

Structure:  
- twitt/pkg -- All the implementations are here for the reuse purpose.  
- twitt/web -- Web side program.
- twitt/server -- Server side program.
