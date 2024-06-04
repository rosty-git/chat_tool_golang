package userusecase

import (
	"log/slog"
	"strings"
	"time"

	"github.com/elef-git/chat_tool_golang/internal/handler"
	"github.com/elef-git/chat_tool_golang/internal/models"
)

type userService interface {
	Registration(userName, email, password string) error
	Login(email, password string) (string, error)
	GetById(userID string) (*models.User, error)
	UpdateStatus(userID string, status string, manual bool, dndEndTime string) (*models.Status, error)
	GetStatus(userID string) (*models.Status, error)
	GetNotUpdatedStatuses() ([]*models.Status, error)
}

type channelService interface {
	GetByUserId(userID string, channelType models.ChannelType) ([]*models.Channel, error)
	GetUsers(channelID string) ([]*models.User, error)
}

type UseCase struct {
	userService      userService
	channelService   channelService
	broadcastChannel chan handler.WsMessage
}

func NewUseCase(userService userService, channelService channelService, broadcastChannel chan handler.WsMessage) *UseCase {
	return &UseCase{
		userService:      userService,
		channelService:   channelService,
		broadcastChannel: broadcastChannel,
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

func (uc *UseCase) GetUsersByChannelId(channelID string) ([]*models.User, error) {
	return uc.channelService.GetUsers(channelID)
}

func (uc *UseCase) UpdateStatus(userID string, status string, manual bool, dndEndTime string) (*models.Status, error) {
	slog.Info("UseCase UpdateStatus", "userID", userID, "status", status, "manual", manual, "dndEndTime", dndEndTime)

	if status == "online" && manual {
		manual = false
	}

	currentStatus, err := uc.userService.GetStatus(userID)
	if err != nil {
		return nil, err
	}

	slog.Info("UseCase UpdateStatus", "currentStatus", currentStatus)

	if currentStatus.Manual && !manual && status != "online" {
		return currentStatus, nil
	}

	newStatus, err := uc.userService.UpdateStatus(userID, status, manual, dndEndTime)
	if err != nil {
		return nil, err
	}

	slog.Info("UseCase UpdateStatus", "newStatus", newStatus)

	if newStatus.Status != currentStatus.Status {
		message := map[string]string{"userId": userID, "status": newStatus.Status}

		wsMessage := handler.WsMessage{Action: "status-updated", Payload: message}

		uc.broadcastChannel <- wsMessage
	}

	return newStatus, nil
}

func (uc *UseCase) GetStatus(userID string) (*models.Status, error) {
	return uc.userService.GetStatus(userID)
}

func (uc *UseCase) StatusesWatchdog() {
	for {
		slog.Info("StatusesWatchdog")

		statuses, err := uc.userService.GetNotUpdatedStatuses()
		if err != nil {
			slog.Error("GetNotUpdatedStatuses", "err", err)
		}

		slog.Info("StatusesWatchdog", "statuses", statuses)

		for _, status := range statuses {
			_, err := uc.UpdateStatus(status.UserID, "offline", false, "")
			if err != nil {
				slog.Error("StatusesWatchdog", "err", err)
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
