package utility

import "github.com/gin-gonic/gin"

func PrepareJsonResponse(status string, statusCode int, data interface{}) gin.H {
	return gin.H{"status": status, "status_code": statusCode, "data": data}
}
