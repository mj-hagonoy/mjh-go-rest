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
	"github.com/mj-hagonoy/mjh-go-rest/pkg/mail"
)

func main() {
	configFile := flag.String("config", "config.yaml", "configuration file")
	flag.Parse()
	if err := config.ParseConfig(*configFile); err != nil {
		panic(err)
	}
	runJobWorker()
	runMailWorker()
	runRestService()
}

func runRestService() {
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.HomeHandler)

	//Users
	router.HandleFunc("/users", handlers.GetAllUsersHandler).Methods(http.MethodGet)
	router.HandleFunc("/users/import", handlers.ImportUsersHandler).Methods(http.MethodPost)

	//Jobs
	router.HandleFunc("/jobs/{id}", handlers.GetJob).Methods(http.MethodGet)

	http.Handle("/", router)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.GetConfig().Port), nil); err != nil {
		log.Fatal(err)
	}
}

func runJobWorker() {
	go func() {
		for {
			req := <-job.JobRequests
			err := job.ProcessJob(context.Background(), req)
			if err != nil {
				fmt.Printf("error processing job %s, %s\n", req.ID, err.Error())
			}
			mail.MailRequests <- mail.Mail{
				Subject:   fmt.Sprintf("[JOB_NOTICE] ID: %s", req.ID),
				EmailTo:   []string{req.InitiatedBy},
				EmailFrom: config.GetConfig().Mail.EmaiFrom,
				Data:      map[string]string{"job_id": req.ID},
				Type:      mail.MAIL_TYPE_JOB_NOTIF,
			}
		}
	}()
}

func runMailWorker() {
	go func() {
		for {
			req := <-mail.MailRequests
			if err := mail.ProcessEmail(&req); err != nil {
				fmt.Printf("error sending email with data = [%+v], error =%s\n", req, err.Error())
			}
		}
	}()
}
