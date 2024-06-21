package service

import (
	"gauravgn90/gin-crud-with-auth/v2/connection"
	"os"
	"strconv"
	"time"
)

// Function to Invalid the token passed from Authorization

func InvalidateToken(token string) bool {
	expiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY"))
	connection.GetRedis().Set(token, true, time.Minute*time.Duration(expiry))
	return true
}

// Function to check if the token is blacklisted
func IsTokenBlackListed(token string) bool {
	// Check if token is blacklisted
	ok, err := connection.GetRedis().Get(token).Result()
	if err != nil {
		return false
	}
	if ok == "" {
		return false
	}
	return true
}
