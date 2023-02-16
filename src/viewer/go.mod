module viewer

go 1.20

require shared/api/nsq v0.0.1

replace shared/api/nsq => ../shared/api/nsq

require (
	github.com/golang/snappy v0.0.1 // indirect
	github.com/nsqio/go-nsq v1.1.0 // indirect
)
