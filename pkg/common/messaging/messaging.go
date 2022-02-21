package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
)

type Message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type IMessagingClient interface {
	Publish(ctx context.Context, msg Message, exchangeName string, exchangeType string) (string, error)
	Close()
}

func GetMessagingClient(projectID string) (IMessagingClient, error) {
	var connErr error
	gcpOnce.Do(func() {
		client, err := pubsub.NewClient(context.Background(), projectID)
		if err != nil {
			connErr = err
			return
		}
		gcpPubSubClient = GcpMessagingClient{
			cl: client,
		}
	})
	return gcpPubSubClient, connErr
}

type GcpMessagingClient struct {
	cl *pubsub.Client
}

var (
	gcpPubSubClient GcpMessagingClient
	gcpOnce         sync.Once
)

func (cl GcpMessagingClient) Close() {
	cl.cl.Close()
}

func (cl GcpMessagingClient) Publish(ctx context.Context, msg Message, exchangeName string, exchangeType string) (string, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %v", err)
	}

	topic := cl.cl.Topic(exchangeName)
	defer topic.Stop()
	res := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	msgID, err := res.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("Publish: %v", err)
	}
	return msgID, nil
}
