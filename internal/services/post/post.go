package postservice

import (
	"github.com/elef-git/chat_tool_golang/internal/handler"
	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type postRepository interface {
	GetByChannelId(db *gorm.DB, channelID string, limit int, before string, after string) ([]*models.Post, error)
	Get(db *gorm.DB, id string) (*models.Post, error)
	Create(db *gorm.DB, userID string, channelID string, message string) (*models.Post, error)
	Search(db *gorm.DB, userID, text string) ([]*models.Post, error)
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

func (s *Service) GetByChannelId(db *gorm.DB, channelID string, limit int, before string, after string) ([]*models.Post, error) {
	return s.postRepository.GetByChannelId(db, channelID, limit, before, after)
}

func (s *Service) GetById(db *gorm.DB, id string) (*models.Post, error) {
	return s.postRepository.Get(db, id)
}

func (s *Service) Create(db *gorm.DB, userID string, channelID string, message string) (*models.Post, error) {
	return s.postRepository.Create(db, userID, channelID, message)
}

func (s *Service) NotifyReceivers(IDs []string, message interface{}) {
	wsMessage := handler.WsMessage{ToUsersIDs: IDs, Action: "new-post", Payload: message}

	s.wsChannel <- wsMessage
}

func (s *Service) NotifySender(userID string, message interface{}) {
	wsMessage := handler.WsMessage{
		ToUsersIDs: []string{userID},
		Action:     "new-own-post",
		Payload:    message,
	}

	s.wsChannel <- wsMessage
}

func (s *Service) Search(db *gorm.DB, userID, text string) ([]*models.Post, error) {
	return s.postRepository.Search(db, userID, text)
}
