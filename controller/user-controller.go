package controller

import (
	"gauravgn90/gin-crud-with-auth/v2/model"
	"gauravgn90/gin-crud-with-auth/v2/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
