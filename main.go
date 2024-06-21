package main

import (
	"fmt"
	"gauravgn90/gin-crud-with-auth/v2/connection"
	db "gauravgn90/gin-crud-with-auth/v2/connection"
	"gauravgn90/gin-crud-with-auth/v2/route"
	"gauravgn90/gin-crud-with-auth/v2/utility"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	_ "gauravgn90/gin-crud-with-auth/v2/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
	// Initialize database connection
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utility.GetEnv("DB_USER"),
		utility.GetEnv("DB_PASSWORD"),
		utility.GetEnv("DB_HOST"),
		utility.GetEnv("DB_PORT"),
		utility.GetEnv("DB_NAME"),
	)
	// Initialize the database connection pool
	if err := db.InitDB(dataSourceName, 100, 10); err != nil {
		log.Fatalf("Error initializing database: %v", err)
		return
	}
	db.RunMigration()
	// Initialize Redis connection
	connection.InitRedis()
}

//	@title			GIN CRUD AND AUTH API
//	@version		1.0
//	@description	This is simple CRUD API with JWT Auth.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8083
//	@BasePath	/api/v1

//	@securityDefinitions.apiKey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				"Enter your bearer token in the format 'Bearer {token}'"
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	defer db.Close()
	router := gin.New()
	route.InitializeRouter(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//run on port 8083
	router.Run(":" + utility.GetEnv("APP_PORT"))
}
