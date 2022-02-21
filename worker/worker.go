package worker

import (
	"fmt"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
)

type IWorker interface {
	Run()
}

type WorkerType string

func GetWorker(workerType WorkerType) (IWorker, error) {
	switch workerType {
	case WORKER_JOB:
		return JobWorker{ProjectID: config.GetConfig().Messaging.GoogleCloud.ProjectID}, nil
	case "web":
		return WebWorker{}, nil
	case WORKER_EMAIL:
		return MailWorker{ProjectID: config.GetConfig().Messaging.GoogleCloud.ProjectID}, nil
	default:
		return nil, fmt.Errorf("main: unsupported type %v", workerType)
	}
}
