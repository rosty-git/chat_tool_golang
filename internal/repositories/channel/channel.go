package channelrepository

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

func (r *Repository) GetByUserId(userID string, channelType models.ChannelType) ([]*models.Channel, error) {
	slog.Info("channelRepo GetByUserId", "userID", userID, "channelType", channelType)

	var channels []*models.Channel

	result := r.db.Table("channels").
		Joins("JOIN channel_members ON channels.id = channel_members.channel_id").
		Where("channel_members.user_id = ? AND channels.type = ? AND channel_members.deleted_at IS NULL", userID, channelType).
		Find(&channels)

	return channels, result.Error
}
