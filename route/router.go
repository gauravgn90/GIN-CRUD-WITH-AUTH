package route

import (
	"gauravgn90/gin-crud-with-auth/v2/controller"
	"gauravgn90/gin-crud-with-auth/v2/middleware"
	"gauravgn90/gin-crud-with-auth/v2/model"
	"gauravgn90/gin-crud-with-auth/v2/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User model.User

func InitializeRouter(router *gin.Engine) *gin.Engine {
	// Set JSON request/response for all APIs
	router.Use(middleware.Logger())
	router.Use(middleware.Cors())
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
	})

	v1 := router.Group("/api/v1")
	{
		usersRoute := v1.Group("/users", middleware.TokenAuthMiddleware())
		{
			usersRoute.GET("/", controller.GetUsers)
			usersRoute.POST("/", controller.CreateUser)
			usersRoute.DELETE("/:id", controller.DeleteUser)
			usersRoute.PUT("/:id", controller.UpdateUser)
		}
		authRoute := v1.Group("/auth")
		{
			authRoute.POST("/login", controller.Login)
			authRoute.POST("/logout", middleware.TokenAuthMiddleware(), controller.Logout)
		}

		v1.POST("/checkToken", middleware.TokenAuthMiddleware(), func(c *gin.Context) {
			c.IndentedJSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, "Token is valid"))
		})
	}

	// fall back to a 404 if the route is not found
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusNotFound, utility.PrepareJsonResponse("error", http.StatusNotFound, "Route not found"))
	})

	return router
}
