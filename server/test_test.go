package main

import (
    "context"
    "fmt"
    "github.com/golang/mock/gomock"
    "github.com/golang/protobuf/proto"
    "sort"
    "testing"
    "time"
    twittmock "twitt/pkg/mock"
    "twitt/pkg/rpc"
    "twitt/pkg/server"
)


var s = server.Server{}

// rpcMsg implements the gomock.Matcher interface
type rpcMsg struct {
    msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
    m, ok := msg.(proto.Message)
    if !ok {
        return false
    }
    return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
    return fmt.Sprintf("is %s", r.msg)
}

// Implement sort.Interface for []*pb.Post.
type postSlice []*pb.Post

func (s postSlice) Len() int { return len(s) }

func (s postSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s postSlice) Less(i, j int) bool {
    if s[i].Username != s[j].Username {
        return s[i].Username < s[j].Username
    }
    return s[i].Content < s[j].Content
}

func TestSignUp(t *testing.T) {
    // Mock TwittServiceClient.
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockTwittServiceClient := twittmock.NewMockTwittServiceClient(ctrl)
    req := &pb.InfoRequest{Username: "test1"}
    mockTwittServiceClient.EXPECT().SignUp(
        gomock.Any(),
        &rpcMsg{msg: req},
    ).Return(&pb.SuccessReply{Success: true}, nil)
    testSignUp(t, mockTwittServiceClient)

    // Set up concrete test cases.
    testcases := []struct{
        username string
        password string
        success bool
    } {
        // Sign up with a new username.
        {
            username: "test1",
            password: "123",
            success: true,
        },
        // Sign up with another username.
        {
            username: "test2",
            password: "111",
            success: true,
        },
        // Sign up with existing username.
        {
            username: "test1",
            password: "321",
            success: false,
        },
        // Sign up with a new username.
        {
            username: "test3",
            password: "aaa",
            success: true,
        },
    }

    for _, testcase := range testcases {
        req := &pb.InfoRequest{Username: testcase.username, Password: testcase.password}
        resp, err := s.SignUp(context.Background(), req)
        if err != nil {
            t.Errorf("TestSignUp got unexpected error")
        }
        if resp.Success != testcase.success {
            t.Errorf("SignUp(%v) got: %v wanted: %v", testcase.username, resp.Success, testcase.success)
        }
    }
}

func testSignUp(t *testing.T, client pb.TwittServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := client.SignUp(ctx, &pb.InfoRequest{Username: "test1"})
    if err != nil || r.Success != true {
        t.Errorf("mocking failed")
    }
    t.Log("Reply : ", r.Success)
}

func TestLogin(t *testing.T) {
    // Mock TwittServiceClient.
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockTwittServiceClient := twittmock.NewMockTwittServiceClient(ctrl)
    req := &pb.InfoRequest{Username: "test1", Password: "222"}
    mockTwittServiceClient.EXPECT().Login(
        gomock.Any(),
        &rpcMsg{msg: req},
    ).Return(&pb.SuccessReply{Success: false}, nil)
    testLogin(t, mockTwittServiceClient)

    // Set up concrete test cases.
    testcases := []struct{
        username string
        password string
        success bool
    } {
        // Log in using non-existing username.
        {
            username: "unknown",
            password: "123",
            success: false,
        },
        // Log in with correct username but wrong password.
        {
            username: "test1",
            password: "error",
            success: false,
        },
        // Log in with correct username and password.
        {
            username: "test1",
            password: "123",
            success: true,
        },
        // Log in with correct username and password.
        {
            username: "test2",
            password: "111",
            success: true,
        },
    }

    for _, testcase := range testcases {
        req := &pb.InfoRequest{Username: testcase.username, Password: testcase.password}
        resp, err := s.Login(context.Background(), req)
        if err != nil {
            t.Errorf("TestLogin got unexpected error")
        }
        if resp.Success != testcase.success {
            t.Errorf("Login(%v) got: %v wanted: %v", testcase.username, resp.Success, testcase.success)
        }
    }
}

func testLogin(t *testing.T, client pb.TwittServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := client.Login(ctx, &pb.InfoRequest{Username: "test1", Password: "222"})
    if err != nil || r.Success != false {
        t.Errorf("mocking failed")
    }
    t.Log("Reply : ", r.Success)
}

func TestPosting(t *testing.T) {
    // Mock TwittServiceClient.
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockTwittServiceClient := twittmock.NewMockTwittServiceClient(ctrl)
    post := &pb.Post{Username: "test2", Content: "This is a post by test2"}
    req := &pb.PostRequest{Post: post}
    mockTwittServiceClient.EXPECT().Posting(
        gomock.Any(),
        &rpcMsg{msg: req},
    ).Return(&pb.SuccessReply{Success: true}, nil)
    testPosting(t, mockTwittServiceClient)

    // Set up concrete test cases.
    testcases := []struct{
        username string
        content string
        success bool
    } {
        // test1 post something.
        {
            username: "test1",
            content: "Wonderful",
            success: true,
        },
        // test3 post something.
        {
            username: "test3",
            content: "Thank you",
            success: true,
        },
        // non-existing user post something.
        {
            username: "error",
            content: "This should not success",
            success: false,
        },
    }

    for _, testcase := range testcases {
        post := &pb.Post{Username: testcase.username, Content: testcase.content}
        req := &pb.PostRequest{Post: post}
        resp, err := s.Posting(context.Background(), req)
        if err != nil {
            t.Errorf("TestPosting got unexpected error")
        }
        if resp.Success != testcase.success {
            t.Errorf("Posting(%v) got: %v wanted: %v", testcase.username, resp.Success, testcase.success)
        }
    }
}

func testPosting(t *testing.T, client pb.TwittServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    post := &pb.Post{Username: "test2", Content: "This is a post by test2"}
    r, err := client.Posting(ctx, &pb.PostRequest{Post: post})
    if err != nil || r.Success != true {
        t.Errorf("mocking failed")
    }
    t.Log("Reply : ", r.Success)
}

func TestFollow(t *testing.T) {
    // Mock TwittServiceClient.
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockTwittServiceClient := twittmock.NewMockTwittServiceClient(ctrl)
    req := &pb.FollowingRequest{Username: "test1", Following: "test3"}
    mockTwittServiceClient.EXPECT().Follow(
        gomock.Any(),
        &rpcMsg{msg: req},
    ).Return(&pb.SuccessReply{Success: true}, nil)
    testFollow(t, mockTwittServiceClient)

    // Set up concrete test cases.
    testcases := []struct{
        username string
        following string
        success bool
    } {
        // test2 follow test1.
        {
            username: "test2",
            following: "test1",
            success: true,
        },
        // test3 follow test1.
        {
            username: "test3",
            following: "test1",
            success: true,
        },
        // test2 follow test3.
        {
            username: "test2",
            following: "test3",
            success: true,
        },
        // test3 follow a non-existing user.
        {
            username: "test3",
            following: "error",
            success: false,
        },
    }

    for _, testcase := range testcases {
        req := &pb.FollowingRequest{Username: testcase.username, Following: testcase.following}
        resp, err := s.Follow(context.Background(), req)
        if err != nil {
            t.Errorf("TestFollow got unexpected error")
        }
        if resp.Success != testcase.success {
            t.Errorf("Follow(%v) got: %v wanted: %v", testcase.username, resp.Success, testcase.success)
        }
    }
}

func testFollow(t *testing.T, client pb.TwittServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := client.Follow(ctx, &pb.FollowingRequest{Username: "test1", Following: "test3"})
    if err != nil || r.Success != true {
        t.Errorf("mocking failed")
    }
    t.Log("Reply : ", r.Success)
}

func TestUnFollow(t *testing.T) {
    // Mock TwittServiceClient.
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockTwittServiceClient := twittmock.NewMockTwittServiceClient(ctrl)
    req := &pb.FollowingRequest{Username: "test1", Following: "test3"}
    mockTwittServiceClient.EXPECT().UnFollow(
        gomock.Any(),
        &rpcMsg{msg: req},
    ).Return(&pb.SuccessReply{Success: true}, nil)
    testUnFollow(t, mockTwittServiceClient)

    // Set up concrete test cases.
    testcases := []struct{
        username string
        following string
        success bool
    } {
        // test2 follow test3 when test3 is following test3.
        {
            username: "test2",
            following: "test3",
            success: true,
        },
        // test3 unfollow test2 when test3 is not following test2.
        {
            username: "test3",
            following: "test2",
            success: false,
        },
        // non-existing user unfollow test1.
        {
            username: "error",
            following: "test",
            success: false,
        },
    }

    for _, testcase := range testcases {
        req := &pb.FollowingRequest{Username: testcase.username, Following: testcase.following}
        resp, err := s.UnFollow(context.Background(), req)
        if err != nil {
            t.Errorf("TestUnFollow got unexpected error")
        }
        if resp.Success != testcase.success {
            t.Errorf("UnFollow(%v) got: %v wanted: %v", testcase.username, resp.Success, testcase.success)
        }
    }
}

func testUnFollow(t *testing.T, client pb.TwittServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := client.UnFollow(ctx, &pb.FollowingRequest{Username: "test1", Following: "test3"})
    if err != nil || r.Success != true {
        t.Errorf("mocking failed")
    }
    t.Log("Reply : ", r.Success)
}

func TestGetList(t *testing.T) {
    // Mock TwittServiceClient.
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockTwittServiceClient := twittmock.NewMockTwittServiceClient(ctrl)
    req := &pb.InfoRequest{Username: "test1", Instruct: "Follow"}
    list := make([]string, 0)
    list = append(list, "test2")
    list = append(list, "test3")
    mockTwittServiceClient.EXPECT().GetList(
        gomock.Any(),
        &rpcMsg{msg: req},
    ).Return(&pb.ListReply{Success: true, List: list}, nil)
    testGetList(t, mockTwittServiceClient)

    // Set up concrete test cases.
    list1 := make([]string, 0)
    list1 = append(list1, "test2")
    list1 = append(list1, "test3")
    list2 := make([]string, 0)
    list2 = append(list2, "test1")
    testcases := []struct{
        username string
        instruct string
        success bool
        list []string
    } {
        // test1 follow noboby.
        {
            username: "test1",
            instruct: "Follow",
            success: true,
            list: list1,
        },
        // test2 follow test1.
        {
            username: "test2",
            instruct: "UnFollow",
            success: true,
            list: list2,
        },
        // instruct is incorrect.
        {
            username: "test2",
            instruct: "error",
            success: false,
            list: nil,
        },
        // user does not exist.
        {
            username: "unknown",
            instruct: "Follow",
            success: false,
            list: nil,
        },
    }

    for _, testcase := range testcases {
        req := &pb.InfoRequest{Username: testcase.username, Instruct: testcase.instruct}
        resp, err := s.GetList(context.Background(), req)
        if err != nil {
            t.Errorf("TestGetList got unexpected error")
        }
        if resp.Success != testcase.success {
            t.Errorf("GetList(%v) got: %v wanted: %v", testcase.username, resp.Success, testcase.success)
        }
        if !stringSliceEqual(resp.List, testcase.list) {
            t.Errorf("GetList(%v) list is not correct", testcase.username)
            for _, p := range resp.List {
                t.Errorf("resp.List: %v", p)
            }
            for _, q := range testcase.list {
                t.Errorf("testcase.list: %v", q)
            }
        }
    }
}

func stringSliceEqual(a, b []string) bool {
    if len(a) != len(b) {
        fmt.Printf("a: %v b: %v", len(a), len(b))
        return false
    }

    if (a == nil) != (b == nil) {
        return false
    }

    sort.Strings(a)
    sort.Strings(b)

    for i, v := range a {
        if v != b[i] {
            fmt.Printf("v: %v b[i]: %v", v, b[i])
            return false
        }
    }

    return true
}

func testGetList(t *testing.T, client pb.TwittServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := client.GetList(ctx, &pb.InfoRequest{Username: "test1", Instruct: "Follow"})
    if err != nil || r.Success != true || r.List[0] != "test2" || r.List[1] != "test3" {
        t.Errorf("mocking failed")
    }
    t.Log("Reply : ", r.Success)
}

func TestView(t *testing.T) {
    // Mock TwittServiceClient.
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockTwittServiceClient := twittmock.NewMockTwittServiceClient(ctrl)
    req := &pb.InfoRequest{Username: "test1"}
    posts := make(postSlice, 0)
    post := &pb.Post{Username: "test1", Content: "Wonderful"}
    posts = append(posts, post)
    mockTwittServiceClient.EXPECT().View(
        gomock.Any(),
        &rpcMsg{msg: req},
    ).Return(&pb.ViewReply{Success: true, Posts: posts}, nil)
    testView(t, mockTwittServiceClient)

    // Set up concrete test cases.
    post1 := &pb.Post{Username: "test1", Content: "Wonderful"}
    post2 := &pb.Post{Username: "test3", Content: "Thank you"}
    posts1 := make(postSlice, 0)
    posts2 := make(postSlice, 0)
    posts3 := make(postSlice, 0)
    posts1 = append(posts1, post1)
    posts2 = append(posts2, post1)
    posts3 = append(posts3, post1)
    posts3 = append(posts3, post2)
    // test2 and test3 both follow test1.
    testcases := []struct{
        username string
        posts postSlice
        success bool
    } {
        {
            username: "test1",
            posts: posts1,
            success: true,
        },
        {
            username: "test2",
            posts: posts2,
            success: true,
        },
        {
            username: "test3",
            posts: posts3,
            success: true,
        },
        {
            username: "unknown",
            posts: nil,
            success: false,
        },
    }

    for _, testcase := range testcases {
        req := &pb.InfoRequest{Username: testcase.username}
        resp, err := s.View(context.Background(), req)
        if err != nil {
            t.Errorf("TestView got unexpected error")
        }
        if resp.Success != testcase.success {
            t.Errorf("View(%v) got: %v wanted: %v", testcase.username, resp.Success, testcase.success)
        }
        if !postSliceEqual(resp.Posts, testcase.posts) {
            t.Errorf("View(%v) post list is not correct", testcase.username)
        }
    }
}

func postSliceEqual(a, b postSlice) bool {
    if len(a) != len(b) {
        return false
    }

    if (a == nil) != (b == nil) {
        return false
    }

    sort.Stable(a)
    sort.Stable(b)

    for i, v := range a {
        if v.Username != b[i].Username || v.Content != b[i].Content {
            return false
        }
    }

    return true
}

func testView(t *testing.T, client pb.TwittServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    r, err := client.View(ctx, &pb.InfoRequest{Username: "test1"})
    if err != nil || r.Success != true || r.Posts[0].Username != "test1" || r.Posts[0].Content != "Wonderful" {
        t.Errorf("mocking failed")
    }
    t.Log("Reply : ", r.Success)
}
