package connection

import (
	"gauravgn90/gin-crud-with-auth/v2/logservice"
	"os"

	"github.com/go-redis/redis"
)

// Initialize Redis connection variables
var (
	rdb        *redis.Client
	REDIS_HOST string
	REDIS_PORT string
)

// Initialize Redis connection
func InitRedis() {
	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PORT = os.Getenv("REDIS_PORT")

	rdb = redis.NewClient(&redis.Options{
		Addr:         REDIS_HOST + ":" + REDIS_PORT,
		PoolSize:     100,
		MinIdleConns: 10,
	})
	logservice.Info("Redis connection initialized %s:%s", REDIS_HOST, REDIS_PORT)

	// check connection
	_, err := rdb.Ping().Result()
	if err != nil {
		logservice.Error("Error initializing Redis connection: %v", err)
		return
	}
	logservice.Info("Redis connection initialized successfully")
}

// Get Redis connection instance
func GetRedis() *redis.Client {
	if rdb != nil {
		return rdb
	}
	return nil
}
