package controller

import (
	"errors"
	"gauravgn90/gin-crud-with-auth/v2/model"
	"gauravgn90/gin-crud-with-auth/v2/utility"
	"net/http"

	db "gauravgn90/gin-crud-with-auth/v2/connection"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login User
func Login(c *gin.Context) {
	errChan := make(chan error)
	go loginUser(c, errChan)
	err := <-errChan
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}

}

// Private function to login user
func loginUser(c *gin.Context, errChan chan<- error) {
	var input model.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		errChan <- err
		return
	}

	var user model.User

	err := db.GetDB().Where("username = ?", input.Username).First(&user).Error

	if err != nil {
		errChan <- err
		return
	}

	if user.Id == 0 {
		errChan <- errors.New("user not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		errChan <- errors.New("invalid password")
		return
	}

	token, err := utility.GenerateToken(c, user.Id, user.Username)
	if err != nil {
		errChan <- errors.New("error generating token")
		return
	}

	c.IndentedJSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, gin.H{"token": token}))
	close(errChan)
}

// Logout User
func Logout(c *gin.Context) {
	// Invalidate JWT Token passed in Authorization header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "Authorization header is required"))
		return
	}
	c.JSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, "Logged out successfully"))
}
