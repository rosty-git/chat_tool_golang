package postrepository

import (
	"log/slog"
	"time"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Get(db *gorm.DB, id string) (*models.Post, error) {
	slog.Info("postRepo Get", "id", id)

	var post models.Post

	result := db.Preload("User").First(&post, "id = ?", id)

	return &post, result.Error
}

func (r *Repository) GetByChannelId(db *gorm.DB, channelID string, limit int, before string, after string) ([]*models.Post, error) {
	slog.Info("postRepo GetByChannelId", "channelID", channelID, "limit", limit, "before", before, "after", after)

	var posts []*models.Post
	var result *gorm.DB

	if before != "" {
		beforeCreatedAt, err := time.Parse(time.RFC3339, before)
		if err != nil {
			slog.Error("Error parsing before", "err", err)

			return nil, err
		}

		result = db.Limit(limit).Model(&models.Post{}).Preload("User").Where(
			"channel_id = ? AND created_at < ?", channelID, beforeCreatedAt,
		).Order("created_at desc").Find(&posts)
	} else if after != "" {
		afterCreatedAt, err := time.Parse(time.RFC3339, after)
		if err != nil {
			slog.Error("Error parsing after", "err", err)

			return nil, err
		}

		result = db.Limit(limit).Model(&models.Post{}).Preload("User").Where(
			"channel_id = ? AND created_at > ?", channelID, afterCreatedAt,
		).Order("created_at desc").Find(&posts)
	} else {
		result = db.Limit(limit).Model(&models.Post{}).Preload("User").Where(
			"channel_id = ?", channelID,
		).Order("created_at desc").Find(&posts)
	}

	slog.Info("Result", "error", result.Error)

	return posts, result.Error
}

func (r *Repository) Create(db *gorm.DB, userID string, channelID string, message string) (*models.Post, error) {
	slog.Info("postRepo Create", "userID", userID, "channelID", channelID, "message", message)

	post := &models.Post{
		UserID:    userID,
		ChannelID: channelID,
		Message:   message,
	}
	result := db.Create(post)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.Get(db, post.ID)
}
