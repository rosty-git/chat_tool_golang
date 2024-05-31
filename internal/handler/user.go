package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"github.com/gin-gonic/gin"
)

type UserV1Handler struct {
	config      config
	userUseCase userUseCase
}

func NewUserV1Handler(config config, userUseCase userUseCase) *UserV1Handler {
	return &UserV1Handler{
		config:      config,
		userUseCase: userUseCase,
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

func (uh *UserV1Handler) UpdateStatus(c *gin.Context) {
	slog.Info("UserV1Handler UpdateStatus")

	type StatusForm struct {
		Status     string `json:"status"`
		Manual     bool   `json:"manual"`
		DNDEndTime string `json:"dnd_end_time"`
	}

	var sf StatusForm
	err := c.BindJSON(&sf)
	if err != nil {
		slog.Error("BindJSON", "err", err)
	}

	authUser, ok := c.Get("user")
	if !ok {
		slog.Error("user not found")
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	user := authUser.(*models.User)

	status, err := uh.userUseCase.UpdateStatus(user.ID, sf.Status, sf.Manual, sf.DNDEndTime)
	if err != nil {
		slog.Error("Updated status", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (uh *UserV1Handler) GetChannelMembers(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		slog.Error("user not found")

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	channelUsers, err := uh.userUseCase.GetUsersByChannelId(c.Param("channelID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channels"})
		return
	}

	var usersIDs []string

	for _, channelUser := range channelUsers {
		if channelUser.ID != user.ID {
			usersIDs = append(usersIDs, channelUser.ID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"members": usersIDs})
}
