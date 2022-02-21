package worker

import (
	"context"
	"encoding/json"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/common/messaging"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/logger"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/mail"
)

const WORKER_EMAIL = "email"

type MailWorker struct {
	cl        *pubsub.Client
	once      sync.Once
	ProjectID string
}

func (w MailWorker) Run() {
	logger.InfoLogger.Println("MailWorker.Run: starting worker")
	if err := w.Connect(w.ProjectID); err != nil {
		panic(err)
	}
	defer w.cl.Close()
	sub := w.cl.Subscription(config.GetConfig().Messaging.GoogleCloud.ProjectID)
	cm := make(chan mail.Mail)
	go func() {
		for {
			mail := <-cm
			if err := mail.ProcessEmail(); err != nil {
				logger.ErrorLogger.Printf("error sending email with error: [%s]\n", err.Error())
			}
		}
	}()

	err := sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		var req messaging.Message
		err := json.Unmarshal(m.Data, &req)
		if err != nil {
			logger.ErrorLogger.Printf("json.Unmarshal: %v\n", err)
			return
		}
		if req.Type != WORKER_EMAIL {
			return
		}
		m.Ack()
		var mail mail.Mail
		if err := json.Unmarshal(req.Data, &mail); err != nil {
			logger.ErrorLogger.Printf("json.Unmarshal: %v\n", err)
			return
		}
		cm <- mail
	})
	if err != nil {
		logger.InfoLogger.Printf("sub.Receive: %v\n", err)
	}
}

func (w *MailWorker) Connect(projectID string) error {
	var connErr error
	w.once.Do(func() {
		client, err := pubsub.NewClient(context.Background(), projectID)
		if err != nil {
			connErr = err
			return
		}
		w.cl = client
	})
	return connErr
}
