package utils

import "net/http"

type AppError struct {
	StatusCode	int
	Message		string
	ErrorCode	string
}

var (
	//400
	ErrInvalidRequest = AppError{http.StatusBadRequest, "INVALID_REQUEST", "ERR_4000"}
	ErrUnauthorized = AppError{http.StatusUnauthorized, "UNAUTHORIZED", "ERR_4010"}
	ErrUserNotFound = AppError{http.StatusNotFound, "USER_NOT_FOUND", "ERR_4040"}
	ErrUserConflict = AppError{http.StatusConflict, "USER_ALREADY_EXIST", "ERR_4090"}

	//500
	ErrInternalServer = AppError{http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "ERR_5000"}
)