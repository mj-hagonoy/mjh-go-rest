package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/common/messaging"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/mail"
)

func (job Job) AddToJobQueue() (string, error) {
	cl, err := messaging.GetMessagingClient(config.GetConfig().Messaging.GoogleCloud.ProjectID)
	if err != nil {
		return "", fmt.Errorf("messaging.GetMessagingClient: %v", err)
	}
	b, _ := json.Marshal(job)
	return cl.Publish(context.Background(), messaging.Message{Type: "job", Data: b}, config.GetConfig().Messaging.GoogleCloud.TopicID, "")
}

func (job Job) Notify() (string, error) {
	cl, err := messaging.GetMessagingClient(config.GetConfig().Messaging.GoogleCloud.ProjectID)
	if err != nil {
		return "", fmt.Errorf("messaging.GetMessagingClient: %v", err)
	}
	msg := mail.Mail{
		Subject:   fmt.Sprintf("[JOB_NOTICE] ID: %s", job.ID),
		EmailTo:   []string{job.InitiatedBy},
		EmailFrom: config.GetConfig().Mail.EmaiFrom,
		Data: map[string]string{
			"job_id": job.ID,
			"url":    fmt.Sprintf("%s/jobs/%s", config.GetConfig().ApiUrl(), job.ID),
		},
		Type: mail.MAIL_TYPE_JOB_NOTIF,
	}

	b, _ := json.Marshal(msg)
	return cl.Publish(context.Background(), messaging.Message{Type: "email", Data: b}, config.GetConfig().Messaging.GoogleCloud.TopicID, "")
}
