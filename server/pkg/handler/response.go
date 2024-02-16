package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type error struct {
	Message string `json:"message"`
}

func newErrorMessage(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, error{message})
}

func newOkMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": string(message),
	})
}
