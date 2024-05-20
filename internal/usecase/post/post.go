package postusecase

import "github.com/elef-git/chat_tool_golang/internal/models"

type postService interface {
	GetByChannelId(channelID string, limit int) ([]*models.Post, error)
	Create(userID string, channelID string, message string) (*models.Post, error)
}

type UseCase struct {
	postService postService
}

func NewUseCase(postService postService) *UseCase {
	return &UseCase{
		postService: postService,
	}
}

func (uc *UseCase) GetByChannelId(channelID string, limit int) ([]*models.Post, error) {
	return uc.postService.GetByChannelId(channelID, limit)
}

func (uc *UseCase) Create(userID string, channelID string, message string) (*models.Post, error) {
	return uc.postService.Create(userID, channelID, message)
}
