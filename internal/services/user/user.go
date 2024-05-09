package userservice

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/elef-git/chat_tool_golang/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type userRepository interface {
	Create(user *models.User) (*models.User, error)
}

type Service struct {
	userRepository userRepository
}

func NewService(userRepository userRepository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) Registration(userName, email, password string) error {
	user := &models.User{
		Name:  userName,
		Email: email,
	}

	// Generate salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	user.Salt = hex.EncodeToString(salt)
	slog.Info("Registering user", "user.Salt", user.Salt)

	// Hash password with salt
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password+user.Salt), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	slog.Info("Registering user", "user.Password", user.Password)

	createdUser, err := s.userRepository.Create(user)
	if err != nil {
		return err
	}

	slog.Info("Created user", "user", createdUser)

	return nil
}
