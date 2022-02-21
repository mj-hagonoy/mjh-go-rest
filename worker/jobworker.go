package worker

import (
	"context"
	"encoding/json"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/common/messaging"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/job"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/logger"
)

const WORKER_JOB = "job"

type JobWorker struct {
	cl        *pubsub.Client
	once      sync.Once
	ProjectID string
}

func (w JobWorker) Run() {
	logger.InfoLogger.Println("JobWorker.Run: starting worker")
	if err := w.Connect(w.ProjectID); err != nil {
		panic(err)
	}
	defer w.cl.Close()
	sub := w.cl.Subscription(config.GetConfig().Messaging.GoogleCloud.ProjectID)
	cm := make(chan job.Job)
	go func() {
		for {
			job := <-cm
			if err := job.ProcessJob(context.Background()); err != nil {
				logger.ErrorLogger.Printf("job.ProcessJob: %v\n", err)
			}

			if _, err := job.Notify(); err != nil {
				logger.ErrorLogger.Printf("job.Notif: %v\n", err)
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
		if req.Type != WORKER_JOB {
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
		logger.InfoLogger.Printf("sub.Receive: %v\n", err)
	}
}

func (w *JobWorker) Connect(projectID string) error {
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
