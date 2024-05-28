package postrepository

import (
	"log/slog"
	"time"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Get(id string) (*models.Post, error) {
	slog.Info("postRepo Get", "id", id)

	var post models.Post

	result := r.db.Preload("User").First(&post, "id = ?", id)

	return &post, result.Error
}

func (r *Repository) GetByChannelId(channelID string, limit int, afterCreatedAt time.Time) ([]*models.Post, error) {
	slog.Info("postRepo GetByChannelId", "channelID", channelID, "limit", limit)

	var posts []*models.Post
	result := r.db.Limit(limit).Model(&models.Post{}).Preload("User").Where(
		"channel_id = ? AND created_at > ?", channelID, afterCreatedAt,
	).Order("created_at desc").Find(&posts)

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
	if result.Error != nil {
		return nil, result.Error
	}

	return r.Get(post.ID)
}
