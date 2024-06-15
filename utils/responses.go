package utils

import "github.com/gin-gonic/gin"

func SuccessResponse(message string, data any) gin.H {
	h := gin.H{}
	h["success"] = true
	if message != "" {
		h["message"] = message
	}
	if data != nil {
		h["data"] = data
	}
	return h
}
func ErrorResponse(err error) gin.H {
	h := gin.H{}
	h["success"] = false
	h["error"] = err.Error()
	return h
}
