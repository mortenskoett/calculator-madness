module calculator

go 1.20

require shared/api/nsq v0.0.1

replace shared/api/nsq => ../shared/api/nsq

require (
	google.golang.org/grpc v1.52.3
	google.golang.org/protobuf v1.28.1
)

require github.com/golang/snappy v0.0.1 // indirect

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/nsqio/go-nsq v1.1.0
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
)
