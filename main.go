package main

import (
	"flag"
	"os"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/logger"
	"github.com/mj-hagonoy/mjh-go-rest/worker"
)

var servType worker.WorkerType

func main() {
	worker, err := worker.GetWorker(servType)
	if err != nil {
		panic(err)
	}
	worker.Run()
}

func init() {
	configFile := flag.String("config", "config.yaml", "configuration file")
	serviceType := flag.String("type", "web", "type = [web, job]")
	flag.Parse()
	if err := config.ParseConfig(*configFile); err != nil {
		panic(err)
	}
	servType = worker.WorkerType(*serviceType)
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.GetConfig().Credentials.GoogleCloud)
	logger.InitLoggers()
}
