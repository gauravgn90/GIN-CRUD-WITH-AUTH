package controller

import (
	"gauravgn90/gin-crud-with-auth/v2/model"
	"gauravgn90/gin-crud-with-auth/v2/utility"
	"net/http"

	db "gauravgn90/gin-crud-with-auth/v2/connection"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login User
func Login(c *gin.Context) {

	var input model.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utility.PrepareJsonResponse("error", http.StatusBadRequest, "Invalid request payload"))
		return
	}

	var user model.User

	err := db.GetDB().Where("username = ?", input.Username).First(&user).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, err.Error()))
		return
	}

	if user.Id == 0 {
		c.JSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "User not found"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "Invalid password"))
		return
	}

	token, err := utility.GenerateToken(c, user.Id, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, "Failed to generate token"))
		return
	}

	c.JSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, gin.H{"token": token}))
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
