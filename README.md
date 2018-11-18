# Stage One:
```sh
$ go build #compile the project.  
$ go test #do the unit test.  
$ ./twitt #run the program.  
```
http://localhost:8080/ is the URL for twitt webpage.  

Structure:  
- "/signup" Sign up here ane alert if the username already exists.  
- "/login"  Log in here.
- "/post"  Post twitter here.
- "/view"  View the followings' and own posts.
- "/follow"  Explore all non-folllowed users and make the folllowing. 
- "/unfollow"  Un-follow any current following users.
- "/logout"  Log out and clean the cookies.  

All the data is stored in global variable uList here. Since every function is very simple, they are all combined in one file and handled by main function.
