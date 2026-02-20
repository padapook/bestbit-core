package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/padapook/bestbit-core/internal/account/model"
	"github.com/padapook/bestbit-core/internal/account/service"
	"github.com/padapook/bestbit-core/internal/utils/auth"
	// "log"
)

type UserController interface {
	Register(c *gin.Context)
	GetProfile(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	LoginByShareToken(c *gin.Context)
	GenerateShareToken(c *gin.Context)
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

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginByShareTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	AccountID string `json:"account_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func toUserResponse(user *model.User) UserResponse {
	return UserResponse{
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

func (ctrl *userController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := ctrl.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	// log.Println("'user",user)

	tokens, err := auth.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": LoginResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			User:         toUserResponse(user),
		},
	})
}

func (ctrl *userController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func (ctrl *userController) LoginByShareToken(c *gin.Context) {
	var req LoginByShareTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := ctrl.userService.LoginByShareToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokens, err := auth.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate session tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in via share token successfully",
		"data": LoginResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			User:         toUserResponse(user),
		},
	})
}

func (ctrl *userController) GenerateShareToken(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		username = c.Query("username")
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}

	user, err := ctrl.userService.GetByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	shareTokenString, err := auth.GenerateShareToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate share token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Share token generated successfully",
		"share_token": shareTokenString,
	})
}
