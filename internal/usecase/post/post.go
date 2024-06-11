package postusecase

import (
	"log/slog"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type postService interface {
	GetByChannelId(db *gorm.DB, channelID string, limit int, before string, after string) ([]*models.Post, error)
	Create(db *gorm.DB, userID string, channelID string, message string) (*models.Post, error)
	NotifyReceivers(usersIDs []string, message interface{})
	NotifySender(userID string, message interface{})
}

type channelService interface {
	GetUsers(db *gorm.DB, channelID string) ([]*models.User, error)
	IncrementTotalMsgCount(db *gorm.DB, channelID string) error
	MarkAsRead(db *gorm.DB, channelID string, userID string) error
}

type UseCase struct {
	postService    postService
	channelService channelService
	db             *gorm.DB
}

func NewUseCase(db *gorm.DB, postService postService, channelService channelService) *UseCase {
	return &UseCase{
		db:             db,
		postService:    postService,
		channelService: channelService,
	}
}

func (uc *UseCase) GetByChannelId(channelID string, limit int, before string, after string) ([]*models.Post, error) {
	return uc.postService.GetByChannelId(uc.db, channelID, limit, before, after)
}

func (uc *UseCase) Create(userID string, channelID string, message string, frontId string) (*models.Post, error) {
	var createdPost *models.Post

	err := uc.db.Transaction(func(tx *gorm.DB) error {
		cp, err := uc.postService.Create(tx, userID, channelID, message)
		if err != nil {
			slog.Error("CreatePost", "err", err)

			return nil
		}
		createdPost = cp

		if err = uc.channelService.IncrementTotalMsgCount(tx, channelID); err != nil {
			return err
		}

		if err = uc.channelService.MarkAsRead(tx, channelID, userID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	channelUsers, err := uc.channelService.GetUsers(uc.db, channelID)
	if err != nil {
		return nil, err
	}

	var usersIDs []string

	for _, channelUser := range channelUsers {
		if channelUser.ID != userID {
			usersIDs = append(usersIDs, channelUser.ID)
		}
	}

	uc.postService.NotifyReceivers(usersIDs, createdPost)

	createdPostMap := map[string]interface{}{
		"createdPost": createdPost,
		"frontId":     frontId,
	}

	uc.postService.NotifySender(userID, createdPostMap)

	return createdPost, nil
}
