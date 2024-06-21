package controller

import (
	"errors"
	"gauravgn90/gin-crud-with-auth/v2/model"
	"gauravgn90/gin-crud-with-auth/v2/service"
	"gauravgn90/gin-crud-with-auth/v2/utility"
	"net/http"
	"strings"

	db "gauravgn90/gin-crud-with-auth/v2/connection"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login User
// Add this in swagger
//
//	@Summary		Login User
//	@Description	Login User
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		model.UserLogin				true	"User Login"
//	@Success		200		{object}	model.UserResponseSuccess	"Success"
//	@Failure		400		{object}	model.UserResponseFailure	"Error"
//	@Failure		500		{object}	model.UserResponseFailure	"Error"
//	@Router			/auth/login [post]
func Login(c *gin.Context) {
	errChan := make(chan error)
	go loginUser(c, errChan)
	err := <-errChan
	if err != nil {
		var customErr *utility.ApiResponseError
		if ok := errors.As(err, &customErr); ok {
			// It's a ApiResponseError, we can access the code and message
			c.IndentedJSON(customErr.Code, utility.PrepareJsonResponse("error", customErr.Code, err.Error()))
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, err))
		return
	}

}

// Logout User
// Add this in swagger

// @Summary		Logout User
// @Description	Logout User
// @Tags			auth
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	model.UserResponseLogoutSuccess					"Success"
// @Failure		400	{object}	model.UserResponseLogoutUnAuthorizedFailure		"Error"
// @Failure		500	{object}	model.UserResponseLogoutInternalServerFailure	"Error"
// @Router			/auth/logout [post]
//
// @Security		BearerAuth
func Logout(c *gin.Context) {
	// Invalidate JWT Token passed in Authorization header
	token := c.GetHeader("Authorization")
	if token == "" || len(token) == 0 {
		c.JSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "Authorization header is required"))
		return
	}

	bearerToken := strings.Split(token, "Bearer ")
	if len(bearerToken) != 2 {
		c.JSON(http.StatusUnauthorized, utility.PrepareJsonResponse("error", http.StatusUnauthorized, "Invalid token"))
		return
	}

	service.InvalidateToken(bearerToken[1])
	c.JSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, "Logged out successfully"))
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
		errChan <- utility.NewApiResponseError(http.StatusNotFound, err.Error())
		return
	}

	if user.Id == 0 {
		errChan <- utility.NewApiResponseError(http.StatusNotFound, "user not found")
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

func createEntity(c *gin.Context, entity interface{}, errChan chan<- error) {
	if err := c.ShouldBindJSON(entity); err != nil {
		errChan <- err
		return
	}

	if err := db.GetDB().Create(entity).Error; err != nil {
		errChan <- &utility.ApiResponseError{Code: http.StatusBadRequest, Message: err.Error()}
		return
	}

	c.IndentedJSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, entity))
	close(errChan)
}

func CreateRole(c *gin.Context) {
	errChan := make(chan error)
	var input model.Role
	go createEntity(c, &input, errChan)
	handleError(c, errChan)
}

func CreatePermission(c *gin.Context) {
	errChan := make(chan error)
	var input model.Permission
	go createEntity(c, &input, errChan)
	handleError(c, errChan)
}

func handleError(c *gin.Context, errChan <-chan error) {
	err := <-errChan
	if err != nil {
		var customErr *utility.ApiResponseError
		if ok := errors.As(err, &customErr); ok {
			// It's a ApiResponseError, we can access the code and message
			c.IndentedJSON(customErr.Code, utility.PrepareJsonResponse("error", customErr.Code, err.Error()))
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}
}

func CreateRolePermission(c *gin.Context) {
	errChan := make(chan error)
	go createRolePermission(c, errChan)
	err := <-errChan
	if err != nil {
		var customErr *utility.ApiResponseError
		if ok := errors.As(err, &customErr); ok {
			// It's a ApiResponseError, we can access the code and message
			c.IndentedJSON(customErr.Code, utility.PrepareJsonResponse("error", customErr.Code, err.Error()))
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}
}

func createRolePermission(c *gin.Context, errChan chan<- error) {
	var input model.CreateRolePermissionRequest
	if inputErr := c.ShouldBindJSON(&input); inputErr != nil {
		errChan <- &utility.ApiResponseError{Code: http.StatusBadRequest, Message: inputErr.Error()}
		return
	}

	role := input.Role
	permissions := input.Permissions

	if createErr := db.GetDB().Create(&role).Error; createErr != nil {
		errChan <- createErr
		return
	}

	assocErr := db.GetDB().Model(&role).Association("Permissions").Append(permissions)

	if assocErr != nil {
		errChan <- assocErr
		return
	}

	c.IndentedJSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, gin.H{"role": role}))
	close(errChan)
}

func GetRolesPermissions(c *gin.Context) {
	errChan := make(chan error)
	go getRolesPermissions(c, errChan)
	err := <-errChan
	if err != nil {
		var customErr *utility.ApiResponseError
		if ok := errors.As(err, &customErr); ok {
			// It's a ApiResponseError, we can access the code and message
			c.IndentedJSON(customErr.Code, utility.PrepareJsonResponse("error", customErr.Code, err.Error()))
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}
}

func getRolesPermissions(c *gin.Context, errChan chan<- error) {
	id := c.Param("id")

	var user model.User
	if err := db.GetDB().Where("id = ?", id).First(&user).Error; err != nil {
		errChan <- &utility.ApiResponseError{Code: http.StatusBadRequest, Message: "user not found"}
		return
	}

	var roles []model.Role
	if err := db.GetDB().Model(&user).Association("Roles").Find(&roles); err != nil {
		errChan <- err
		return
	}

	var permissions []model.Permission
	for _, role := range roles {
		var rolePermissions []model.Permission
		if err := db.GetDB().Model(&role).Association("Permissions").Find(&rolePermissions); err != nil {
			errChan <- err
			return
		}
		permissions = append(permissions, rolePermissions...)
	}

	c.IndentedJSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, gin.H{"permissions": permissions, "roles": roles}))
	close(errChan)
}

func AssignRolesToUser(c *gin.Context) {
	errChan := make(chan error)
	go assignRolesToUser(c, errChan)
	err := <-errChan
	if err != nil {
		var customErr *utility.ApiResponseError
		if ok := errors.As(err, &customErr); ok {
			// It's a ApiResponseError, we can access the code and message
			c.IndentedJSON(customErr.Code, utility.PrepareJsonResponse("error", customErr.Code, err.Error()))
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}
}

func assignRolesToUser(c *gin.Context, errChan chan<- error) {
	var input model.AssignRolesToUserRequest
	if inputErr := c.ShouldBindJSON(&input); inputErr != nil {
		errChan <- &utility.ApiResponseError{Code: http.StatusBadRequest, Message: inputErr.Error()}
		return
	}

	user := input.UserID
	roles := input.Roles

	// Check if user exists
	var userModel model.User
	if err := db.GetDB().Where("id = ?", user).First(&userModel).Error; err != nil {
		errChan <- &utility.ApiResponseError{Code: http.StatusBadRequest, Message: "user not found"}
		return
	}

	// Check if roles exist
	var roleModels []model.Role
	if err := db.GetDB().Where("id IN (?)", roles).Find(&roleModels).Error; err != nil {
		errChan <- &utility.ApiResponseError{Code: http.StatusBadRequest, Message: "roles not found"}
		return
	}

	assocErr := db.GetDB().Model(&userModel).Association("Roles").Append(roleModels)

	if assocErr != nil {
		errChan <- assocErr
		return
	}

	c.IndentedJSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, gin.H{"user": user}))
	close(errChan)
}

func AssignPermissionsToRole(c *gin.Context) {
	errChan := make(chan error)
	go assignPermissionsToRole(c, errChan)
	err := <-errChan
	if err != nil {
		var customErr *utility.ApiResponseError
		if ok := errors.As(err, &customErr); ok {
			// It's a ApiResponseError, we can access the code and message
			c.IndentedJSON(customErr.Code, utility.PrepareJsonResponse("error", customErr.Code, err.Error()))
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, utility.PrepareJsonResponse("error", http.StatusInternalServerError, err.Error()))
		return
	}
}

func assignPermissionsToRole(c *gin.Context, errChan chan<- error) {
	var input model.AssignPermissionsToRoleRequest
	if inputErr := c.ShouldBindJSON(&input); inputErr != nil {
		errChan <- &utility.ApiResponseError{Code: http.StatusBadRequest, Message: inputErr.Error()}
		return
	}

	role := input.RoleID
	// Permissions array of permission ids
	permissions := input.Permissions

	// Check if role exists
	var roleModel model.Role
	if err := db.GetDB().Where("id = ?", role).First(&roleModel).Error; err != nil {
		errChan <- &utility.ApiResponseError{Code: http.StatusBadRequest, Message: "role not found"}
		return
	}

	// Check if permissions exist
	var permissionModels []model.Permission
	if err := db.GetDB().Where("id IN (?)", permissions).Find(&permissionModels).Error; err != nil {
		errChan <- &utility.ApiResponseError{Code: http.StatusBadRequest, Message: "permissions not found"}
		return
	}

	assocErr := db.GetDB().Model(&roleModel).Association("Permissions").Append(permissionModels)

	if assocErr != nil {
		errChan <- assocErr
		return
	}

	c.IndentedJSON(http.StatusOK, utility.PrepareJsonResponse("success", http.StatusOK, gin.H{"role": role}))
	close(errChan)
}
