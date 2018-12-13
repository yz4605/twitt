package main

import (
	"flag"
	"go.etcd.io/etcd/raft/raftpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"twitt/pkg/rpc"
	"twitt/pkg/server"
)

func main() {
	cluster := flag.String("cluster", "http://127.0.0.1:9001", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	join := flag.Bool("join", false, "join an existing cluster")
	port := flag.String("port", "2233", "port for grpc")
	flag.Parse()

	proposeC := make(chan string)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)
	server.Start(cluster,id,join,proposeC,confChangeC)

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTwittServiceServer(s, &server.Server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
