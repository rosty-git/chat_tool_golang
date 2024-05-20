package postrepository

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

func (r *Repository) GetByChannelId(channelID string, limit int) ([]*models.Post, error) {
	slog.Info("postRepo GetByChannelId", "channelID", channelID, "limit", limit)

	var posts []*models.Post
	result := r.db.Limit(limit).Model(&models.Post{}).Preload("User").Where(&models.Post{ChannelID: channelID}).Order("created_at asc").Find(&posts)

	slog.Info("Result", "error", result.Error)

	return posts, result.Error
}

func (r *Repository) Create(userID string, channelID string, message string) (*models.Post, error) {
	slog.Info("postRepo Create", "userID", userID, "channelID", channelID, "message", message)

	post := &models.Post{
		UserID:    userID,
		ChannelID: channelID,
		Message:   message,
	}
	result := r.db.Create(post)

	slog.Info("Created post: ", post)

	return post, result.Error
}
