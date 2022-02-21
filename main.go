package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mj-hagonoy/mjh-go-rest/handlers"
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
		runRestService()
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

func runRestService() {
	logger.InfoLogger.Println("HTTP service is adding handlers.")
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.HomeHandler)

	api1 := router.PathPrefix("/api/v1").Subrouter()

	//Users
	api1.HandleFunc("/users", handlers.GetAllUsersHandler).Methods(http.MethodGet)
	api1.HandleFunc("/users/import", handlers.ImportUsersHandler).Methods(http.MethodPost)

	//Jobs
	api1.HandleFunc("/jobs/{id}", handlers.GetJob).Methods(http.MethodGet)

	http.Handle("/", router)

	logger.InfoLogger.Printf("HTTP service will listen to %s at port %d", config.GetConfig().Host, config.GetConfig().Port)
	logger.InfoLogger.Println("HTTP service running...")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.GetConfig().Port), nil); err != nil {
		log.Fatal(err)
	}
}
