package rest

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status int
	Data   interface{}
	Err    error
}

func (r Response) GetJson() []byte {
	var data []byte
	var err error
	if r.Err != nil {
		data, err = json.Marshal(map[string]string{
			"error": r.Err.Error(),
		})
	} else {
		data, err = json.Marshal(r.Data)
	}
	if err != nil {
		return nil
	}
	return data
}

func Respond(w http.ResponseWriter, response Response) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	w.Write(response.GetJson())
}
