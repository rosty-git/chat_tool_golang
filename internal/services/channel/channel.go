package channelservice

import (
	"log/slog"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type channelRepository interface {
	GetByUserId(db *gorm.DB, userId string, channelType models.ChannelType) ([]*models.Channel, error)
	GetMembers(db *gorm.DB, channelID string) ([]*models.ChannelMembers, error)
	GetUsers(db *gorm.DB, channelID string) ([]*models.User, error)
}

type Service struct {
	channelRepository channelRepository
}

func NewService(channelRepository channelRepository) *Service {
	return &Service{
		channelRepository: channelRepository,
	}
}

func (s *Service) GetByUserId(db *gorm.DB, userID string, channelType models.ChannelType) ([]*models.Channel, error) {
	slog.Info("Channel service GetByUserId", "userID", userID, "channelType", channelType)

	return s.channelRepository.GetByUserId(db, userID, channelType)
}

func (s *Service) GetUsers(db *gorm.DB, channelID string) ([]*models.User, error) {
	return s.channelRepository.GetUsers(db, channelID)
}
