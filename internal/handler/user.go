package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type userUseCase interface {
	Registration(userName, email, password string) error
}

type UserV1Handler struct {
	userUseCase userUseCase
}

func NewUserV1Handler(userUseCase userUseCase) *UserV1Handler {
	return &UserV1Handler{
		userUseCase: userUseCase,
	}
}

func (uh *UserV1Handler) Login(c *gin.Context) {
	slog.Info("Login")

	c.JSON(200, gin.H{
		"hello": "login",
	})
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
