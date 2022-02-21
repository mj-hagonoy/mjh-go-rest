package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/logger"
)

var servType string

func main() {
	switch servType {
	case WORKER_JOB:
		worker := JobWorker{ProjectID: config.GetConfig().Messaging.GoogleCloud.ProjectID}
		worker.Run()
	case "web":
		worker := WebWorker{}
		worker.Run()
	case WORKER_EMAIL:
		worker := MailWorker{ProjectID: config.GetConfig().Messaging.GoogleCloud.ProjectID}
		worker.Run()
	default:
		panic(fmt.Sprintf("main: unsupported type %v", servType))
	}
}

func init() {
	configFile := flag.String("config", "config.yaml", "configuration file")
	serviceType := flag.String("type", "web", "type = [web, job]")
	flag.Parse()
	if err := config.ParseConfig(*configFile); err != nil {
		panic(err)
	}
	servType = *serviceType
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.GetConfig().Credentials.GoogleCloud)
	logger.InitLoggers()
}
