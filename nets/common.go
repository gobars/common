package nets

import "net/http"

type Result struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(data interface{}) *Result {
	return &Result{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    data,
	}
}
