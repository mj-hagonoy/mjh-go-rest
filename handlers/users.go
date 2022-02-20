package handlers

import (
	"net/http"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/file_storage"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/job"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/rest"
	"github.com/mj-hagonoy/mjh-go-rest/pkg/user"
)

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := user.GetAll(r.Context())
	if err != nil {
		rest.Respond(w, rest.Response{Status: http.StatusInternalServerError, Err: err})
		return
	}
	if len(users) == 0 {
		rest.Respond(w, rest.Response{Status: http.StatusNoContent})
		return
	}
	rest.Respond(w, rest.Response{Status: http.StatusOK, Data: users})
}

// ImportUsersHandler accepts csv file attachment and stores contents to the users collection
//
// Request Form-data:
//
// file : csv file
func ImportUsersHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	_, header, err := r.FormFile("file")
	if err != nil {
		rest.Respond(w, rest.Response{Status: http.StatusBadRequest, Err: err})
		return
	}
	file, err := header.Open()
	if err != nil {
		rest.Respond(w, rest.Response{Status: http.StatusBadRequest, Err: err})
		return
	}
	var fileData []byte = make([]byte, header.Size)
	if _, err := file.Read(fileData); err != nil {
		rest.Respond(w, rest.Response{Status: http.StatusBadRequest, Err: err})
		return
	}

	storage, err := file_storage.GetStorage()
	if err != nil {
		rest.Respond(w, rest.Response{Status: http.StatusInternalServerError, Err: err})
		return
	}
	if err := storage.Write(r.Context(), header.Filename, fileData); err != nil {
		rest.Respond(w, rest.Response{Status: http.StatusBadRequest, Err: err})
		return
	}

	importJob, err := job.CreateNewJob(r.Context(), job.JOB_TYPE_IMPORT_USERS, header.Filename)
	if err != nil {
		rest.Respond(w, rest.Response{Status: http.StatusInternalServerError, Err: err})
		return
	}
	rest.Respond(w, rest.Response{Status: http.StatusOK, Data: importJob})
}
