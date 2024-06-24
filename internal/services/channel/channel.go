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
	IncrementTotalMsgCount(db *gorm.DB, channelID string) error
	MarkAsRead(db *gorm.DB, channelID string, userID string) error
	GetUnreadCount(db *gorm.DB, channelID string, userID string) (uint64, error)
	SearchOpenChannels(db *gorm.DB, text string) ([]*models.Channel, error)
	GetDirectByMembers(db *gorm.DB, memberID1 string, memberID2 string) (*models.Channel, error)
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

func (s *Service) IncrementTotalMsgCount(db *gorm.DB, channelID string) error {
	return s.channelRepository.IncrementTotalMsgCount(db, channelID)
}

func (s *Service) MarkAsRead(db *gorm.DB, channelID string, userID string) error {
	return s.channelRepository.MarkAsRead(db, channelID, userID)
}

func (s *Service) GetUnreadCount(db *gorm.DB, channelID string, userID string) (uint64, error) {
	return s.channelRepository.GetUnreadCount(db, channelID, userID)
}

func (s *Service) Search(db *gorm.DB, text string) ([]*models.Channel, error) {
	return s.channelRepository.SearchOpenChannels(db, text)
}

func (s *Service) GetDirectByMembers(db *gorm.DB, memberID1 string, memberID2 string) (*models.Channel, error) {
	return s.channelRepository.GetDirectByMembers(db, memberID1, memberID2)
}
