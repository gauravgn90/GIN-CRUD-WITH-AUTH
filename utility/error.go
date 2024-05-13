package utility

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error recovery function
func ErrorRecovery(c *gin.Context) {
	if r := recover(); r != nil {
		fmt.Println("Recovered from panic:", r)
		c.IndentedJSON(http.StatusInternalServerError, PrepareJsonResponse("error", http.StatusInternalServerError, "Internal server error"))
	}
}
