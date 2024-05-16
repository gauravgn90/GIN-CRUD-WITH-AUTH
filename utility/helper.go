package utility

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// ApiResponseError type includes an error code and a message
type ApiResponseError struct {
	Code    int
	Message string
}

// Implement the Error method for the ApiResponseError type
func (e *ApiResponseError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// Function to create a new ApiResponseError
func NewApiResponseError(code int, message string) *ApiResponseError {
	return &ApiResponseError{
		Code:    code,
		Message: message,
	}
}

var jwtKey []byte

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	jwt.StandardClaims
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GenerateToken(c *gin.Context, Id int, Username string) (token string, err error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   Id,
		Username: Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return token, nil
}
