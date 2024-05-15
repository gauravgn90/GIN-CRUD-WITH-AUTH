package main

import (
	"fmt"
	db "gauravgn90/gin-crud-with-auth/v2/connection"
	"gauravgn90/gin-crud-with-auth/v2/route"
	"gauravgn90/gin-crud-with-auth/v2/utility"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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
}

func main() {
	defer db.Close()
	router := gin.New()
	route.InitializeRouter(router)
	//run on port 8083
	router.Run(":" + utility.GetEnv("APP_PORT"))
}
