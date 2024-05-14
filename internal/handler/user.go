package handler

import (
	"errors"
	"github.com/elef-git/chat_tool_golang/internal/models"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type userUseCase interface {
	Registration(userName, email, password string) error
	Login(email, password string) (string, error)
	GetChannels(userId uint64) ([]models.Channel, error)
	GetContacts(userId uint64) ([]models.User, error)
}

type UserV1Handler struct {
	userUseCase userUseCase
	config      config
}

func NewUserV1Handler(userUseCase userUseCase, config config) *UserV1Handler {
	return &UserV1Handler{
		userUseCase: userUseCase,
		config:      config,
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
	}

	channels, err := uh.userUseCase.GetChannels(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channels"})
	}

	c.JSON(http.StatusOK, gin.H{"channels": channels})
}

func (uh *UserV1Handler) GetContacts(c *gin.Context) {
	slog.Info("UserV1Handler GetContacts")

	user, err := getUser(c)
	if err != nil {
		slog.Error("user not found")

		c.JSON(http.StatusInternalServerError, gin.H{})
	}

	contacts, err := uh.userUseCase.GetContacts(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get contacts"})
	}

	c.JSON(http.StatusOK, gin.H{"contacts": contacts})
}
