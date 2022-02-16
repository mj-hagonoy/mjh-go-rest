package handlers

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/mj-hagonoy/mjh-go-rest/pkg/config"
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

	filepath := path.Clean(fmt.Sprintf("%s/%s", config.GetConfig().Directory.UploadUsers, header.Filename))
	ioutil.WriteFile(filepath, fileData, fs.ModePerm)

	importJob, err := job.CreateNewJob(r.Context(), job.JOB_TYPE_IMPORT_USERS, filepath)
	if err != nil {
		rest.Respond(w, rest.Response{Status: http.StatusInternalServerError, Err: err})
		return
	}
	rest.Respond(w, rest.Response{Status: http.StatusOK, Data: importJob})
}
