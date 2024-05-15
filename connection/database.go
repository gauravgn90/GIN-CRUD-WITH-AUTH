package connection

import (
	"gauravgn90/gin-crud-with-auth/v2/model"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

func InitDB(dataSourceName string, maxOpenConns, maxIdleConns int) error {
	// Open a new database connection with a connection pool
	conn, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	// Set maximum open and idle connections

	conn.DB().SetMaxOpenConns(maxOpenConns)
	conn.DB().SetMaxIdleConns(maxIdleConns)

	// Test the connection
	if err := conn.DB().Ping(); err != nil {
		return err
	}

	db = conn

	return nil
}

func GetDB() *gorm.DB {
	// Return the existing database connection
	if db != nil {
		return db
	}
	log.Println("Database connection not found")
	return nil
}

func Close() error {
	// Close the database connection
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}

func RunMigration() {
	// Migrate the schema
	if err := GetDB().AutoMigrate(&model.User{}).Error; err != nil {
		// Handle the error
	}
}
