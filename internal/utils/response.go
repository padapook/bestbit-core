package utils

import (
	"github.com/gin-gonic/gin"
	"errors"
	"gorm.io/gorm"
)

type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	ErrorCode  string      `json:"error_code,omitempty"`
}

func HandleError(c *gin.Context, appErr AppError) {
	c.JSON(appErr.StatusCode, Response{
		Success:   false,
		Message:   appErr.Message,
		ErrorCode: appErr.ErrorCode,
	})
}

func HandleSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func HandleServiceError(c *gin.Context, err error) {
    var appErr AppError
    if errors.As(err, &appErr) {
        HandleError(c, appErr)
        return
    }

	if errors.Is(err, gorm.ErrRecordNotFound) {
		HandleError(c, ErrUserNotFound)
		return
	}
	
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		HandleError(c, ErrUserConflict)
		return
	}

    HandleError(c, ErrInternalServer)
}