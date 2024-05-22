package handler

import (
	"errors"
	"github.com/elef-git/chat_tool_golang/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type UserV1Handler struct {
	config      config
	userUseCase userUseCase
	postUseCase postUseCase
}

func NewUserV1Handler(config config, userUseCase userUseCase, postUseCase postUseCase) *UserV1Handler {
	return &UserV1Handler{
		config:      config,
		userUseCase: userUseCase,
		postUseCase: postUseCase,
	}
}

func (uh *UserV1Handler) Login(c *gin.Context) {
	slog.Info("Login")

	type LoginForm struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var form LoginForm
	err := c.BindJSON(&form)
	if err != nil {
		slog.Error("BindJSON", "err", err)
	}

	slog.Info("Login", "form", form)

	tokenString, err := uh.userUseCase.Login(form.Email, form.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		uh.config.GetAuthCookieName(),
		tokenString,
		uh.config.GetAuthCookieMaxAge(),
		uh.config.GetAuthCookiePath(),
		uh.config.GetAuthCookieDomain(),
		uh.config.GetAuthCookieSecure(),
		uh.config.GetAuthCookieHttpOnly(),
	)

	c.JSON(http.StatusOK, gin.H{})
}

func (uh *UserV1Handler) Registration(c *gin.Context) {
	slog.Info("UserV1Handler Registration")

	type RegForm struct {
		UserName string `json:"userName"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var form RegForm
	err := c.BindJSON(&form)
	if err != nil {
		slog.Error("BindJSON", "err", err)
	}

	slog.Info("Registration", "form", form)

	if err = uh.userUseCase.Registration(form.UserName, form.Email, form.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func getUser(c *gin.Context) (*models.User, error) {
	authUser, ok := c.Get("user")
	if !ok {
		return nil, errors.New("user not found")
	}

	user := authUser.(*models.User)
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (uh *UserV1Handler) GetChannels(c *gin.Context) {
	slog.Info("UserV1Handler GetChannels")

	user, err := getUser(c)
	if err != nil {
		slog.Error("user not found")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	queryChannelType := c.Query("channelType")
	channelType, ok := models.ChannelTypesMap[queryChannelType]
	if !ok || queryChannelType == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "channelType is required"})
		return
	}

	channels, err := uh.userUseCase.GetChannelsByUserId(user.ID, channelType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channels"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"channels": channels})
}
