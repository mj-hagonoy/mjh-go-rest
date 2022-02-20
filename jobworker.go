package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/job"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/logger"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/mail"
)

var pubsubClient *pubsub.Client
var once sync.Once

func runJobWorker() {
	if err := connect(config.GetConfig().Messaging.GoogleCloud.ProjectID); err != nil {
		panic(err)
	}
	defer pubsubClient.Close()

	sub := pubsubClient.Subscription(config.GetConfig().Messaging.GoogleCloud.ProjectID)
	cm := make(chan *pubsub.Message)
	go func() {
		for {
			m := <-cm
			var req job.Job
			if err := json.Unmarshal(m.Data, &req); err != nil {
				logger.ErrorLogger.Printf("json.Unmarshal: %v\n", err)
				return
			}
			if err := job.ProcessJob(context.Background(), req); err != nil {
				logger.ErrorLogger.Printf("job.ProcessJob: %v\n", err)
			}
			mail.MailRequests <- mail.Mail{
				Subject:   fmt.Sprintf("[JOB_NOTICE] ID: %s", req.ID),
				EmailTo:   []string{req.InitiatedBy},
				EmailFrom: config.GetConfig().Mail.EmaiFrom,
				Data: map[string]string{
					"job_id": req.ID,
					"url":    fmt.Sprintf("%s/jobs/%s", config.GetConfig().ApiUrl(), req.ID),
				},
				Type: mail.MAIL_TYPE_JOB_NOTIF,
			}
		}
	}()

	err := sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		cm <- m
		m.Ack()
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
