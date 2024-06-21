package connection

import (
	"gauravgn90/gin-crud-with-auth/v2/logservice"
	"gauravgn90/gin-crud-with-auth/v2/model"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func InitDB(dataSourceName string, maxOpenConns, maxIdleConns int) error {
	// Open a new database connection with a connection pool
	conn, connectionErr := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if connectionErr != nil {
		return connectionErr
	}

	// Set maximum open and idle connections

	// conn.DB().SetMaxOpenConns(maxOpenConns)
	// conn.DB().SetMaxIdleConns(maxIdleConns)

	var sqlDB, _ = conn.DB()
	// Test the connection
	if dbInstanceErr := sqlDB.Ping(); dbInstanceErr != nil {
		return dbInstanceErr
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)

	db = conn

	return nil
}

func GetDB() *gorm.DB {
	// Return the existing database connection
	if db != nil {
		return db
	}
	log.Fatalln("Database connection not found")
	return nil
}

func Close() error {
	// Close the database connection
	if db != nil {
		sqlDB, _ := db.DB()
		return sqlDB.Close()
	}
	return nil
}

func RunMigration() {
	// Migrate the schema
	db := GetDB()
	if err := db.AutoMigrate(&model.User{}); err != nil {
		// Handle the error
		logservice.Error("Error migrating user table: %v", err)
	}

	if err := db.AutoMigrate(&model.Role{}); err != nil {
		// Handle the error
		logservice.Error("Error migrating role table: %v", err)
	}

	if err := db.AutoMigrate(&model.Permission{}); err != nil {
		// Handle the error
		logservice.Error("Error migrating permission table: %v", err)
	}

	if err := db.AutoMigrate(&model.RolePermission{}); err != nil {
		// Handle the error
		logservice.Error("Error migrating role permission table: %v", err)
	}

	if err := db.AutoMigrate(&model.UserRole{}); err != nil {
		// Handle the error
		logservice.Error("Error migrating user role table: %v", err)
	}

	if err := db.AutoMigrate(&model.Product{}); err != nil {
		logservice.Error("Error migrating product table: %v", err)
	}

}
