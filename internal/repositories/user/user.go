package userrepository

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
	"log/slog"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(user *models.User) (*models.User, error) {
	result := r.db.Create(user)

	slog.Info("Created user: ", user)

	return user, result.Error
}

func (r *Repository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)

	return &user, result.Error
}

func (r *Repository) GetById(id string) (*models.User, error) {
	slog.Info("GetById", "id", id)

	var user models.User
	result := r.db.First(&user, "id = ?", id)

	return &user, result.Error
}
