package userusecase

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
	"log/slog"
	"strings"
)

type userService interface {
	Registration(userName, email, password string) error
	Login(email, password string) (string, error)
	GetById(userID string) (*models.User, error)
}

type channelService interface {
	GetByUserId(userID string, channelType models.ChannelType) ([]*models.Channel, error)
}

type postService interface {
	GetByChannelId(channelID string, limit int) ([]*models.Post, error)
}

type UseCase struct {
	userService    userService
	channelService channelService
	postService    postService
}

func NewUseCase(userService userService, channelService channelService, postService postService) *UseCase {
	return &UseCase{
		userService:    userService,
		channelService: channelService,
		postService:    postService,
	}
}

func (uc *UseCase) Registration(userName, email, password string) error {
	return uc.userService.Registration(userName, email, password)
}

func (uc *UseCase) Login(email, password string) (string, error) {
	return uc.userService.Login(email, password)
}

func (uc *UseCase) GetChannelsByUserId(userID string, channelType models.ChannelType) ([]*models.Channel, error) {
	channels, err := uc.channelService.GetByUserId(userID, channelType)
	if err != nil {
		return nil, err
	}

	if channelType == models.ChannelTypeDirect {
		for _, channel := range channels {
			contactID := strings.Replace(channel.Name, userID, "", 1)
			contactID = strings.Replace(contactID, "__", "", 1)

			slog.Info("GetChannelsByUserId", "contactID", contactID)

			user, err := uc.userService.GetById(contactID)
			if err != nil {
				return nil, err
			}

			channel.Name = user.Name

			slog.Info("GetChannelsByUserId", "channel", channel)
		}
	}

	return channels, err
}
