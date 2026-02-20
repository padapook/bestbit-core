package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/padapook/bestbit-core/internal/account/model"
	"github.com/padapook/bestbit-core/internal/account/service"
)

type UserController interface {
	Register(c *gin.Context)
	GetProfile(c *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{userService: userService}
}

type UserRegisterRequest struct {
	Username  string `json:"username" binding:"required,min=4"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type UserResponse struct {
	UID       uint64 `json:"uid"`
	AccountID string `json:"account_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func toUserResponse(user *model.User) UserResponse {
	return UserResponse{
		UID:       user.ID,
		AccountID: user.AccountId,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func (ctrl *userController) Register(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userModel := model.User{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	createdUser, err := ctrl.userService.Register(c.Request.Context(), userModel)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    toUserResponse(createdUser),
	})
}

func (ctrl *userController) GetProfile(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	user, err := ctrl.userService.GetByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": toUserResponse(user),
	})
}
