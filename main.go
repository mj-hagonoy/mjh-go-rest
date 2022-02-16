package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mj-hagonoy/mjh-go-rest/handlers"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/job"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/logger"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/mail"
)

func main() {
	configFile := flag.String("config", "config.yaml", "configuration file")
	flag.Parse()
	if err := config.ParseConfig(*configFile); err != nil {
		panic(err)
	}
	logger.InitLoggers()
	runJobWorker()
	runMailWorker()
	runRestService()
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

func runJobWorker() {
	logger.InfoLogger.Println("Job batch worker started.")
	go func() {
		for {
			req := <-job.JobRequests
			err := job.ProcessJob(context.Background(), req)
			if err != nil {
				logger.ErrorLogger.Printf("error processing job [%s] with error:[%s]\n", req.ID, err.Error())
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
