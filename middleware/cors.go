package middleware

import (
	"gauravgn90/gin-crud-with-auth/v2/logservice"
	"gauravgn90/gin-crud-with-auth/v2/utility"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		mp := make(map[string]string)
		mp["origin"] = origin
		mp["url"] = utility.GetEnv("FRONTEND_URL")
		logservice.Info("Origin: %s", mp)
		var url = utility.GetEnv("FRONTEND_URL")
		if origin != url {
			log.Printf("Origin not allowed: %s", origin)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "Origin not allowed"))
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", url)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatusJSON(http.StatusNoContent, utility.PrepareJsonResponse("success", http.StatusNoContent, "Preflight request"))
			return
		}

		c.Next()
	}
}

// Path: middleware/jwt-auth.go
