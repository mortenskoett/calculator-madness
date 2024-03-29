module viewer

go 1.20

replace shared/queue => ../shared/queue

require (
	google.golang.org/grpc v1.55.0
	google.golang.org/protobuf v1.30.0
	shared/queue v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
)

require (
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/nsqio/go-nsq v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)
