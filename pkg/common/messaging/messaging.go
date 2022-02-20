package messaging

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
)

type IMessagingClient interface {
	Publish(ctx context.Context, msg []byte, exchangeName string, exchangeType string) (string, error)
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

func (cl GcpMessagingClient) Publish(ctx context.Context, msg []byte, exchangeName string, exchangeType string) (string, error) {
	topic := cl.cl.Topic(exchangeName)
	defer topic.Stop()
	res := topic.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	msgID, err := res.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("Publish: %v", err)
	}
	return msgID, nil
}
