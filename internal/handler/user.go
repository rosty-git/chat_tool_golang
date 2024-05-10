package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type userUseCase interface {
	Registration(userName, email, password string) error
	Login(email, password string) (string, error)
}

type config interface {
	GetEnv() string
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
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", true, true)

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
