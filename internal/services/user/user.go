package userservice

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/elef-git/chat_tool_golang/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type userRepository interface {
	Create(user *models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetById(Id string) (*models.User, error)
	CreateOrUpdateStatus(userID string, status string, manual bool, dndEndTime time.Time) (*models.Status, error)
}

type config interface {
	GetJwtSecret() string
	GetJwtTtl() time.Duration
}

type Service struct {
	userRepository userRepository
	config         config
}

func NewService(userRepository userRepository, config config) *Service {
	return &Service{
		userRepository: userRepository,
		config:         config,
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

func (s *Service) Login(email, password string) (string, error) {
	user, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return "", err
	}

	// Compare hashed passwords
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+user.Salt)); err != nil {
		slog.Error("Wrong password", "err", err)

		return "", err
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(s.config.GetJwtTtl()).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(s.config.GetJwtSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) GetById(UserID string) (*models.User, error) {
	return s.userRepository.GetById(UserID)
}

func (s *Service) UpdateStatus(userID string, status string, manual bool, endTime time.Time) (*models.Status, error) {
	return s.userRepository.CreateOrUpdateStatus(userID, status, manual, endTime)
}
