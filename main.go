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
	"github.com/mj-hagonoy/mjh-go-rest/pkg/mail"
)

var servType string

func main() {
	runMailWorker()
	switch servType {
	case "job":
		runJobWorker()
	case "web":
		runRestService()
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

func runMailWorker() {
	logger.InfoLogger.Println("Mail worker started.")
	go func() {
		for {
			req := <-mail.MailRequests
			if err := mail.ProcessEmail(&req); err != nil {
				logger.ErrorLogger.Printf("error sending email with error: [%s]\n", err.Error())
			}
		}
	}()
}
