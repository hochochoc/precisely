package utils

import (
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	Code   int         `json:"code"`
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

func JsonRespond(w http.ResponseWriter, status bool, code int, err error, data interface{}) {

	response := HttpResponse{
		Status: status,
		Code:   code,
		Data:   data,
	}
	if err != nil {
		response.Error = err.Error()
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
