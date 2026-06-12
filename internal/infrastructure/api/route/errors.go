package apiroute

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidationErrorResponse(err error) gin.H {
	return gin.H{
		"error":   "validation failed",
		"details": err.Error(),
	}
}
