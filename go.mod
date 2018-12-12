module twitt

replace go.etcd.io/etcd => ../src/go.etcd.io/etcd

require (
	github.com/golang/mock v1.1.1
	github.com/golang/protobuf v1.2.0
	go.etcd.io/etcd v3.3.10+incompatible
	go.uber.org/zap v1.9.1
	google.golang.org/grpc v1.16.0
)
