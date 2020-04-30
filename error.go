package gyan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Status  string
	Message string `json:"message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v: %v", r.Status, r.Message)
}

func newResponseError(res *http.Response) error {
	er := &ErrorResponse{Status: res.Status}
	if err := json.NewDecoder(res.Body).Decode(er); err != nil {
		er.Message = err.Error()
	}
	return er
}
