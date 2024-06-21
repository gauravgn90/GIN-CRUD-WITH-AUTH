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
	//router.Use(middleware.Cors())
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
	})

	v1 := router.Group("/api/v1")
	{
		usersRoute := v1.Group("/users" /* middleware.TokenAuthMiddleware() */)
		{
			usersRoute.GET("/", controller.GetUsers)
			usersRoute.POST("/", controller.CreateUser)
			usersRoute.DELETE("/:id", controller.DeleteUser)
			usersRoute.PUT("/:id", controller.UpdateUser)

			requiredRoles := []string{"admin"}
			requiredPermissions := []string{"create"}

			usersRoute.POST("/roles", middleware.AuthenticateUser(requiredRoles, requiredPermissions), controller.CreateRole)
			usersRoute.POST("/permissions", middleware.AuthenticateUser(requiredRoles, requiredPermissions), controller.CreatePermission)

			usersRoute.POST("/roles-permissions", middleware.AuthenticateUser(requiredRoles, requiredPermissions), controller.CreateRolePermission)
			usersRoute.POST("/assign-permissions-to-role", middleware.AuthenticateUser(requiredRoles, requiredPermissions), controller.AssignPermissionsToRole)
			usersRoute.POST("/assign-roles-to-user", middleware.AuthenticateUser(requiredRoles, requiredPermissions), controller.AssignRolesToUser)

			// routes to list all roles and permissions of the user
			usersRoute.GET("/roles-permissions/:id", controller.GetRolesPermissions)
		}
		authRoute := v1.Group("/auth")
		{
			authRoute.POST("/login", controller.Login)
			authRoute.POST("/logout", middleware.TokenAuthMiddleware(), controller.Logout)
		}

		/* productRoute := v1.Group("/products")
		{
			productRoute.GET("/", controller.GetProducts)
			productRoute.POST("/", controller.CreateProduct)
			productRoute.DELETE("/:id", controller.DeleteProduct)
			productRoute.PUT("/:id", controller.UpdateProduct)
		}
		*/
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
