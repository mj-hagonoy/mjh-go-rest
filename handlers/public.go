package handlers

import (
	"net/http"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/rest"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	rest.Respond(w, rest.Response{
		Status: http.StatusOK,
		Data: map[string]string{
			"message": "OK",
		},
	})
}
