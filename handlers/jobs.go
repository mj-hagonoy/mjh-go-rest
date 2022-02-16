package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/job"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/rest"
)

func GetJob(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		rest.Respond(w, rest.Response{Status: http.StatusBadRequest, Err: fmt.Errorf("missing variable id")})
		return
	}
	data, err := job.NewJob(job.ID(id))
	if err != nil {
		rest.Respond(w, rest.Response{Status: http.StatusInternalServerError, Err: err})
		return
	}
	err = data.GetOne(r.Context())
	if err != nil {
		rest.Respond(w, rest.Response{Status: http.StatusInternalServerError, Err: err})
		return
	}
	rest.Respond(w, rest.Response{Status: http.StatusOK, Data: data})
}
