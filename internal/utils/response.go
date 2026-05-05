package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, errMsg string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   errMsg,
	})
}

func BadRequestResponse(c *gin.Context, errMsg string) {
	ErrorResponse(c, http.StatusBadRequest, errMsg)
}

func UnauthorizedResponse(c *gin.Context, errMsg string) {
	ErrorResponse(c, http.StatusUnauthorized, errMsg)
}

func ForbiddenResponse(c *gin.Context, errMsg string) {
	ErrorResponse(c, http.StatusForbidden, errMsg)
}

func NotFoundResponse(c *gin.Context, errMsg string) {
	ErrorResponse(c, http.StatusNotFound, errMsg)
}

func InternalServerErrorResponse(c *gin.Context, errMsg string) {
	ErrorResponse(c, http.StatusInternalServerError, errMsg)
}