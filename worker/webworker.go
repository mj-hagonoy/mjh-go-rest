package worker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mj-hagonoy/mjh-go-rest/handlers"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/logger"
)

type WebWorker struct{}

func (w WebWorker) Run() {
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
