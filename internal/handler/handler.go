package handler

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
)

type config interface {
	GetEnv() string
	GetAuthCookieName() string
	GetAuthCookieMaxAge() int
	GetAuthCookiePath() string
	GetAuthCookieDomain() string
	GetAuthCookieSecure() bool
	GetAuthCookieHttpOnly() bool
	GetCorsAllowOrigins() []string
}

type userUseCase interface {
	Registration(userName, email, password string) error
	Login(email, password string) (string, error)
	GetChannelsByUserId(userId string, channelType models.ChannelType) ([]*models.Channel, error)
	UpdateStatus(userID string, status string, manual bool, dndEndTime string) (*models.Status, error)
}

type postUseCase interface {
	GetByChannelId(channelID string, limit int, before string, after string) ([]*models.Post, error)
	Create(userID string, channelID string, message string) (*models.Post, error)
}
