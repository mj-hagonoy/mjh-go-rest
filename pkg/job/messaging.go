package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/common/messaging"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
)

func AddToQueue(job Job) (string, error) {
	cl, err := messaging.GetMessagingClient(config.GetConfig().Messaging.GoogleCloud.ProjectID)
	if err != nil {
		return "", fmt.Errorf("messaging.GetMessagingClient: %v", err)
	}
	b, _ := json.Marshal(job)
	return cl.Publish(context.Background(), b, config.GetConfig().Messaging.GoogleCloud.TopicID, "")
}
