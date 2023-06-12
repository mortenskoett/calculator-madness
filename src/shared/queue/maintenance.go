package queue

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type queueMaintainer struct {
	nsqlookupAddr string
	client        http.Client
}

func NewQueueMaintainer(nsqHttpAddr string) *queueMaintainer {
	return &queueMaintainer{
		nsqlookupAddr: nsqHttpAddr,
		client:        http.Client{Timeout: 5 * time.Second},
	}
}

func (m *queueMaintainer) DeleteTopic(topic string) error {
	log.Println("requesting deletion of topic", topic)

	requrl := fmt.Sprintf("http://%s/topic/delete?topic=%s", m.nsqlookupAddr, topic)
	resp, err := m.client.Post(requrl, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to make http request to delete topic: %v\n", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read message body: %v", err)
	}

	if resp != nil {
		log.Printf("receieved topic deletion response: %d %s\n", resp.StatusCode, string(b))
	}
	return nil
}

func (m *queueMaintainer) CreateTopic(topic string) error {
	log.Println("requesting creation of topic", topic)

	requrl := fmt.Sprintf("http://%s/topic/create?topic=%s", m.nsqlookupAddr, topic)
	resp, err := m.client.Post(requrl, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to make http request to create topic: %v\n", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read message body: %v", err)
	}

	if resp != nil {
		log.Printf("receieved topic creation response: %d %s\n", resp.StatusCode, string(b))
	}

	return m.createChannel(topic, topic+"-channel")
}

func (m *queueMaintainer) createChannel(topic, channel string) error {
	log.Println("requesting creation of channel", channel, "for topic", topic)

	requrl := fmt.Sprintf("http://%s/channel/create?topic=%s&channel=%s", m.nsqlookupAddr, topic, channel)
	resp, err := m.client.Post(requrl, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to make http request to create channel: %v\n", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read message body: %v", err)
	}

	if resp != nil {
		log.Printf("receieved channel creation response: %d %s\n", resp.StatusCode, string(b))
	}
	return nil
}
