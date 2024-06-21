package utility

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// ApiResponseError type includes an error code and a message
type ApiResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
	Username string   `json:"username"`
	UserID   int      `json:"user_id"`
	Roles    []string `json:"roles"`
	jwt.StandardClaims
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GenerateToken(c *gin.Context, Id int, Username string) (token string, err error) {
	expireInHour, _ := strconv.Atoi(os.Getenv("TOKEN_EXPIRE_IN_MS"))
	expirationTime := time.Now().Add(time.Duration(expireInHour) * time.Millisecond)
	claims := &Claims{
		UserID:   Id,
		Username: Username,
		Roles:    []string{"admin", "editor"},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Audience:  "http://localhost:8083",
			Issuer:    "http://localhost:8083",
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Subject:   "auth",
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return token, nil
}
