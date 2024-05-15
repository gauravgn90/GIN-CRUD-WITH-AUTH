package middleware

import (
	"gauravgn90/gin-crud-with-auth/v2/utility"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		authToken := strings.Split(bearerToken, "Bearer ")

		if len(authToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "Unauthorized access"))
			return
		}
		JWT_SECRET := utility.GetEnv("JWT_SECRET")

		token, err := jwt.Parse(authToken[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(JWT_SECRET), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "Unauthorized access"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "Unauthorized access"))
			return
		}

		userID := int(claims["user_id"].(float64))
		// Optionally, you can set the user ID in the context for further use in handlers
		c.Set("userID", userID)
		c.Next()
	}
}
