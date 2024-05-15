package controller

import (
	"errors"
	"fmt"
	"gauravgn90/gin-crud-with-auth/v2/model"
	"gauravgn90/gin-crud-with-auth/v2/service"
	"gauravgn90/gin-crud-with-auth/v2/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	userService    service.UserService = service.New()
	userController UserController      = New(userService)
)

type UserController interface {
	SaveUser(ctx *gin.Context) (model.User, int, error)
	FindAll() ([]model.User, error)
	Delete(ctx *gin.Context) error
	Update(ctx *gin.Context) error
}

type UserControllerImpl struct {
	service service.UserService
}

// FindAll implements UserController.
func (u UserControllerImpl) FindAll() ([]model.User, error) {
	users, err := u.service.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// SaveUser implements UserController.
func (u UserControllerImpl) SaveUser(ctx *gin.Context) (model.User, int, error) {
	var user model.User
	if err := ctx.BindJSON(&user); err != nil {
		return model.User{}, http.StatusBadRequest, err
	}

	user, err := u.service.SaveUser(user)
	var customErr *utility.CustomError

	if ok := errors.As(err, &customErr); ok {
		// It's a CustomError, we can access the code and message
		return model.User{}, customErr.Code, customErr
	}

	if err != nil {
		return model.User{}, http.StatusInternalServerError, err
	}
	return user, http.StatusCreated, nil
}

// Delete user
func (u UserControllerImpl) Delete(ctx *gin.Context) error {
	// convert string to int
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return err
	}
	delErr := u.service.Delete(id)
	if delErr != nil {
		return delErr
	}
	return nil
}

// Update User
func (u UserControllerImpl) Update(ctx *gin.Context) error {
	// convert string to int
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return err
	}
	var user model.User
	if err := ctx.BindJSON(&user); err != nil {
		return err
	}

	updateErr := u.service.Update(id, user)
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func New(service service.UserService) UserController {
	return UserControllerImpl{
		service: service,
	}
}

// Get Users
func GetUsers(c *gin.Context) {
	fmt.Printf("User Id passed via middleware : %d\n", c.MustGet("userID"))
	defer func() {
		if r := recover(); r != nil {
			c.IndentedJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, "Internal Server Error"))
		}
	}()
	users, err := userController.FindAll()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}
	var userResponses []model.UserResponse
	for _, user := range users {
		userResponse := model.UserResponse{
			Id:       user.Id,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		}
		userResponses = append(userResponses, userResponse)
	}

	c.IndentedJSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, userResponses))
}

// Create User
func CreateUser(c *gin.Context) {
	user, statusCode, err := userController.SaveUser(c)
	if err != nil {
		c.IndentedJSON(statusCode, utility.PrepareJsonResponse("error", statusCode, err.Error()))
		return
	}
	response := model.UserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	}
	c.IndentedJSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, response))
}

// Delete User
func DeleteUser(c *gin.Context) {
	err := userController.Delete(c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}
	c.IndentedJSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, "User deleted successfully"))
}

// Update User
func UpdateUser(c *gin.Context) {
	err := userController.Update(c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}
	c.IndentedJSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, "User updated successfully"))
}
