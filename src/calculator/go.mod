module calculator

go 1.20

replace shared/queue => ../shared/queue

require (
	google.golang.org/grpc v1.52.3
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/nsqio/go-nsq v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
	shared/queue v0.0.0-00010101000000-000000000000
)
