package postusecase

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
	"log/slog"
	"time"
)

type postService interface {
	GetByChannelId(channelID string, limit int, afterCreatedAt time.Time) ([]*models.Post, error)
	Create(userID string, channelID string, message string) (*models.Post, error)
	NotifyReceivers(userID []string, message interface{})
}

type channelService interface {
	GetUsers(channelID string) ([]*models.User, error)
}

type UseCase struct {
	postService    postService
	channelService channelService
}

func NewUseCase(postService postService, channelService channelService) *UseCase {
	return &UseCase{
		postService:    postService,
		channelService: channelService,
	}
}

func (uc *UseCase) GetByChannelId(channelID string, limit int, afterCreatedAt time.Time) ([]*models.Post, error) {
	return uc.postService.GetByChannelId(channelID, limit, afterCreatedAt)
}

func (uc *UseCase) Create(userID string, channelID string, message string) (*models.Post, error) {
	createdPost, err := uc.postService.Create(userID, channelID, message)
	if err != nil {
		slog.Error("CreatePost", "err", err)

		return nil, err
	}

	channelUsers, err := uc.channelService.GetUsers(channelID)
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
