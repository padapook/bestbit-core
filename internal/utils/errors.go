package utils

import "net/http"

type AppError struct {
	StatusCode	int
	Message		string
	ErrorCode	string
}

var (
	ErrInvalidRequest = AppError{http.StatusBadRequest, "Invalid Request", "ERR_4000"}
)