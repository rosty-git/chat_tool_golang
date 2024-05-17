package channelservice

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
	"log/slog"
)

type channelRepository interface {
	GetByUserId(userId string, channelType models.ChannelType) ([]*models.Channel, error)
}

type Service struct {
	channelRepository channelRepository
}

func NewService(channelRepository channelRepository) *Service {
	return &Service{
		channelRepository: channelRepository,
	}
}

func (s Service) GetByUserId(userID string, channelType models.ChannelType) ([]*models.Channel, error) {
	slog.Info("Channel service GetByUserId", "userID", userID, "channelType", channelType)

	return s.channelRepository.GetByUserId(userID, channelType)
}
