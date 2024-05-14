package userusecase

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
)

type userService interface {
	Registration(userName, email, password string) error
	Login(email, password string) (string, error)
	GetChannels(userID uint64) ([]models.Channel, error)
	GetContacts(userId uint64) ([]models.User, error)
}

type UseCase struct {
	userService userService
}

func NewUseCase(userService userService) *UseCase {
	return &UseCase{
		userService: userService,
	}
}

func (uc *UseCase) Registration(userName, email, password string) error {
	return uc.userService.Registration(userName, email, password)
}

func (uc *UseCase) Login(email, password string) (string, error) {
	return uc.userService.Login(email, password)
}

func (uc *UseCase) GetChannels(userID uint64) ([]models.Channel, error) {
	return uc.userService.GetChannels(userID)
}

func (uc *UseCase) GetContacts(userID uint64) ([]models.User, error) {
	return uc.userService.GetContacts(userID)
}
