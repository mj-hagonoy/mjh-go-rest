package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/common/messaging"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/job"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/logger"
)

var pubsubClient *pubsub.Client
var once sync.Once

func runJobWorker() {
	if err := connect(config.GetConfig().Messaging.GoogleCloud.ProjectID); err != nil {
		panic(err)
	}
	defer pubsubClient.Close()

	sub := pubsubClient.Subscription(config.GetConfig().Messaging.GoogleCloud.ProjectID)
	cm := make(chan job.Job)
	go func() {
		for {
			job := <-cm
			if err := job.ProcessJob(context.Background()); err != nil {
				logger.ErrorLogger.Printf("job.ProcessJob: %v\n", err)
			}

			job.Notify()
		}
	}()

	err := sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		var req messaging.Message
		err := json.Unmarshal(m.Data, &req)
		if err != nil {
			logger.ErrorLogger.Printf("json.Unmarshal: %v\n", err)
			return
		}
		if req.Type != "job" {
			return
		}
		m.Ack()
		var job job.Job
		if err := json.Unmarshal(req.Data, &job); err != nil {
			logger.ErrorLogger.Printf("json.Unmarshal: %v\n", err)
			return
		}
		cm <- job
	})
	if err != nil {
		fmt.Printf("sub.Receive: %v\n", err)
	}
}

func connect(projectID string) error {
	var connErr error
	once.Do(func() {
		client, err := pubsub.NewClient(context.Background(), projectID)
		if err != nil {
			connErr = err
			return
		}
		pubsubClient = client
	})
	return connErr
}
