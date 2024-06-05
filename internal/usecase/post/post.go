package postusecase

import (
	"log/slog"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type postService interface {
	GetByChannelId(db *gorm.DB, channelID string, limit int, before string, after string) ([]*models.Post, error)
	Create(db *gorm.DB, userID string, channelID string, message string) (*models.Post, error)
	NotifyReceivers(userID []string, message interface{})
}

type channelService interface {
	GetUsers(db *gorm.DB, channelID string) ([]*models.User, error)
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

func (uc *UseCase) Create(userID string, channelID string, message string) (*models.Post, error) {
	createdPost, err := uc.postService.Create(uc.db, userID, channelID, message)
	if err != nil {
		slog.Error("CreatePost", "err", err)

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

	return createdPost, nil
}
