package postservice

import (
	"github.com/elef-git/chat_tool_golang/internal/handler"
	"github.com/elef-git/chat_tool_golang/internal/models"
)

type postRepository interface {
	GetByChannelId(channelID string, limit int) ([]*models.Post, error)
	Create(userID string, channelID string, message string) (*models.Post, error)
}

type Service struct {
	postRepository postRepository
	wsChannel      chan handler.WsMessage
}

func NewService(channelRepository postRepository, wsChannel chan handler.WsMessage) *Service {
	return &Service{
		postRepository: channelRepository,
		wsChannel:      wsChannel,
	}
}

func (s *Service) GetByChannelId(channelID string, limit int) ([]*models.Post, error) {
	return s.postRepository.GetByChannelId(channelID, limit)
}

func (s *Service) Create(userID string, channelID string, message string) (*models.Post, error) {
	return s.postRepository.Create(userID, channelID, message)
}

func (s *Service) NotifyReceivers(IDs []string, message interface{}) {
	wsMessage := handler.WsMessage{ToUsersIDs: IDs, Action: "new-post", Payload: message}

	s.wsChannel <- wsMessage
}
