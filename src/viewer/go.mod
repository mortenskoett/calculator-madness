module viewer

go 1.20

replace shared/queue => ../shared/queue

require shared/queue v0.0.0-00010101000000-000000000000

require (
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.0
	github.com/nsqio/go-nsq v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)
