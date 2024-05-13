package main

import (
	"gauravgn90/gin-crud-with-auth/v2/connection"
	"gauravgn90/gin-crud-with-auth/v2/controller"
	"gauravgn90/gin-crud-with-auth/v2/model"
	"gauravgn90/gin-crud-with-auth/v2/service"
	. "gauravgn90/gin-crud-with-auth/v2/utility"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User model.User

var (
	userService    service.UserService       = service.New()
	userController controller.UserController = controller.New(userService)
)

func init() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
	// Initialize database connection
	dataSourceName := "root:root@tcp(localhost:3306)/golang"
	// Initialize the database connection pool
	if err := connection.InitDB(dataSourceName, 100, 10); err != nil {
		log.Fatalf("Error initializing database: %v", err)
		return
	}
}

func main() {
	defer connection.GetDB().Close()
	router := gin.Default()
	// Set JSON request/response for all APIs
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
	})

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", getUsers)
		v1.POST("/users", createUser)
		v1.DELETE("/users/:id", deleteUser)
		v1.PUT("/users/:id", updateUser)
	}

	// fall back to a 404 if the route is not found
	router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusNotFound, PrepareJsonResponse("error", http.StatusNotFound, "Route not found"))
	})
	//run on port 8083
	router.Run(":8083")
}

// Get Users
func getUsers(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.IndentedJSON(http.StatusInternalServerError, PrepareJsonResponse("error", http.StatusInternalServerError, "Internal Server Error"))
		}
	}()
	users, err := userController.FindAll()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}
	c.IndentedJSON(http.StatusOK, PrepareJsonResponse("success", http.StatusOK, users))
}

// Create User
func createUser(c *gin.Context) {
	user, statusCode, err := userController.SaveUser(c)
	if err != nil {
		c.IndentedJSON(statusCode, PrepareJsonResponse("error", statusCode, err.Error()))
		return
	}
	user.Password = "**********"
	c.IndentedJSON(http.StatusOK, PrepareJsonResponse("success", http.StatusOK, user))
}

// Delete User
func deleteUser(c *gin.Context) {
	err := userController.Delete(c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}
	c.IndentedJSON(http.StatusOK, PrepareJsonResponse("success", http.StatusOK, "User deleted successfully"))
}

// Update User
func updateUser(c *gin.Context) {
	err := userController.Update(c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}
	c.IndentedJSON(http.StatusOK, PrepareJsonResponse("success", http.StatusOK, "User updated successfully"))
}
