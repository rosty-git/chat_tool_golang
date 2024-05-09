package userrepository

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
	"log/slog"
)

type Repository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(user *models.User) (*models.User, error) {
	result := r.db.Create(user)
	slog.Info("Result", "Error", result.Error, "RowsAffected", result.RowsAffected)

	slog.Info("Created user: ", user)

	return user, nil
}
