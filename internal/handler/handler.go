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
	GetAwsRegion() string
	GetAwsS3Bucket() string
}

type userUseCase interface {
	Registration(userName, email, password string) error
	Login(email, password string) (string, error)
	GetChannelsByUserId(userId string, channelType models.ChannelType) ([]*models.Channel, error)
	UpdateStatus(userID string, status string, manual bool, dndEndTime string) (*models.Status, error)
	GetStatus(userID string) (*models.Status, error)
	GetUsersByChannelId(channelID string) ([]*models.User, error)
}

type channelUseCase interface {
	MarkAsRead(channelID string, userID string) error
	GetUnreadCount(channelID string, userID string) (uint64, error)
}

type postUseCase interface {
	GetByChannelId(channelID string, limit int, before string, after string) ([]*models.Post, error)
	Create(userID string, channelID string, message string, frontId string, files []string) (*models.Post, error)
	Search(userID string, text string) ([]*models.Post, error)
}

type fileUseCase interface {
	CreateTmp(name, fileType string, size uint64) (*models.TmpFile, error)
	SetS3Key(id, s3Key string) (*models.TmpFile, error)
	DeleteTmp(id string) error
	GetPresignUrl(key string) (string, error)
}
