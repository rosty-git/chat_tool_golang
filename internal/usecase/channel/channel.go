package channelusecase

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type channelService interface {
	GetUsers(db *gorm.DB, channelID string) ([]*models.User, error)
	IncrementTotalMsgCount(db *gorm.DB, channelID string) error
	MarkAsRead(db *gorm.DB, channelID string, userID string) error
	GetUnreadCount(db *gorm.DB, channelID string, userID string) (uint64, error)
}

type UseCase struct {
	channelService channelService
	db             *gorm.DB
}

func NewUseCase(db *gorm.DB, channelService channelService) *UseCase {
	return &UseCase{
		db:             db,
		channelService: channelService,
	}
}

func (uc *UseCase) MarkAsRead(channelID string, userID string) error {
	return uc.channelService.MarkAsRead(uc.db, channelID, userID)
}

func (uc *UseCase) GetUnreadCount(channelID string, userID string) (uint64, error) {
	return uc.channelService.GetUnreadCount(uc.db, channelID, userID)
}
