package postservice

import "github.com/elef-git/chat_tool_golang/internal/models"

type postRepository interface {
	GetByChannelId(channelID string, limit int) ([]*models.Post, error)
	Create(userID string, channelID string, message string) (*models.Post, error)
}

type Service struct {
	postRepository postRepository
}

func NewService(channelRepository postRepository) *Service {
	return &Service{
		postRepository: channelRepository,
	}
}

func (s *Service) GetByChannelId(channelID string, limit int) ([]*models.Post, error) {
	return s.postRepository.GetByChannelId(channelID, limit)
}

func (s *Service) Create(userID string, channelID string, message string) (*models.Post, error) {
	return s.postRepository.Create(userID, channelID, message)
}
