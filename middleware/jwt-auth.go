package middleware

import (
	"gauravgn90/gin-crud-with-auth/v2/service"
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

		var isBlacklisted = service.IsTokenBlackListed(authToken[1])
		if isBlacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "Token Blacklisted"))
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
		c.Set("roles", claims["roles"])
		c.Next()
	}
}

func AuthenticateUser(requiredRoles []string, requiredPermissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		roles, _ := c.Get("roles")

		if userID == nil || roles == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "Unauthorized access: user ID or roles not defined"))
			return
		}

		// Check if the user has one of the required roles
		isAuthorized := false
		rolesSlice, ok := roles.([]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, "Internal Server Error"))
			return
		}

		for _, userRole := range rolesSlice {
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					isAuthorized = true
					break
				}
			}
			if isAuthorized {
				break
			}
		}

		if !isAuthorized {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "Unauthorized access: user does not have the required roles"))
			return
		}

		c.Next()
	}
}
