package userusecase

import (
	"log/slog"
	"strings"
	"time"

	"github.com/elef-git/chat_tool_golang/internal/handler"
	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type userService interface {
	Registration(db *gorm.DB, userName, email, password string) error
	Login(db *gorm.DB, email, password string) (string, error)
	GetById(db *gorm.DB, userID string) (*models.User, error)
	UpdateStatus(db *gorm.DB, userID string, status string, manual bool, dndEndTime string) (*models.Status, error)
	GetStatus(db *gorm.DB, userID string) (*models.Status, error)
	GetNotUpdatedStatuses(db *gorm.DB) ([]*models.Status, error)
}

type channelService interface {
	GetByUserId(db *gorm.DB, userID string, channelType models.ChannelType) ([]*models.Channel, error)
	GetUsers(db *gorm.DB, channelID string) ([]*models.User, error)
}

type UseCase struct {
	userService      userService
	channelService   channelService
	broadcastChannel chan handler.WsMessage
	db               *gorm.DB
}

func NewUseCase(db *gorm.DB, userService userService, channelService channelService, broadcastChannel chan handler.WsMessage) *UseCase {
	return &UseCase{
		db:               db,
		userService:      userService,
		channelService:   channelService,
		broadcastChannel: broadcastChannel,
	}
}

func (uc *UseCase) Registration(userName, email, password string) error {
	return uc.userService.Registration(uc.db, userName, email, password)
}

func (uc *UseCase) Login(email, password string) (string, error) {
	return uc.userService.Login(uc.db, email, password)
}

func (uc *UseCase) GetChannelsByUserId(userID string, channelType models.ChannelType) ([]*models.Channel, error) {
	channels, err := uc.channelService.GetByUserId(uc.db, userID, channelType)
	if err != nil {
		return nil, err
	}

	if channelType == models.ChannelTypeDirect {
		for _, channel := range channels {
			contactID := strings.Replace(channel.Name, userID, "", 1)
			contactID = strings.Replace(contactID, "__", "", 1)

			slog.Info("GetChannelsByUserId", "contactID", contactID)

			user, err := uc.userService.GetById(uc.db, contactID)
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
	return uc.channelService.GetUsers(uc.db, channelID)
}

func (uc *UseCase) UpdateStatus(userID string, status string, manual bool, dndEndTime string) (*models.Status, error) {
	if status == "online" && manual {
		manual = false
	}

	currentStatus, err := uc.userService.GetStatus(uc.db, userID)
	if err != nil {
		return nil, err
	}

	if currentStatus.Manual && !manual && status != "online" {
		return currentStatus, nil
	}

	newStatus, err := uc.userService.UpdateStatus(uc.db, userID, status, manual, dndEndTime)
	if err != nil {
		return nil, err
	}

	if newStatus.Status != currentStatus.Status {
		message := map[string]string{"userId": userID, "status": newStatus.Status}

		wsMessage := handler.WsMessage{Action: "status-updated", Payload: message}

		uc.broadcastChannel <- wsMessage
	}

	return newStatus, nil
}

func (uc *UseCase) GetStatus(userID string) (*models.Status, error) {
	return uc.userService.GetStatus(uc.db, userID)
}

func (uc *UseCase) StatusesWatchdog() {
	for {
		statuses, err := uc.userService.GetNotUpdatedStatuses(uc.db)
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

		time.Sleep(5 * time.Minute)
	}
}
