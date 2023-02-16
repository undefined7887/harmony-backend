grpc:
	protoc internal/pkg/centrifugo/proto/centrifugo.proto --go_out=internal/pkg/centrifugo/proto --go-grpc_out=internal/pkg/centrifugo/proto
