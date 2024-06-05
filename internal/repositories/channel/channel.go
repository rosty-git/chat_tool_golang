package channelrepository

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"

	"log/slog"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetByUserId(db *gorm.DB, userID string, channelType models.ChannelType) ([]*models.Channel, error) {
	slog.Info("channelRepo GetByUserId", "userID", userID, "channelType", channelType)

	var channels []*models.Channel

	result := db.Table("channels").
		Joins("JOIN channel_members ON channels.id = channel_members.channel_id").
		Where("channel_members.user_id = ? AND channels.type = ? AND channel_members.deleted_at IS NULL", userID, channelType).
		Find(&channels)

	return channels, result.Error
}

func (r *Repository) GetMembers(db *gorm.DB, channelID string) ([]*models.ChannelMembers, error) {
	var channelMembers []*models.ChannelMembers

	result := db.Preload("User").Find(&channelMembers, "channel_id = ?", channelID)
	if result.Error != nil {
		return nil, result.Error
	}

	return channelMembers, nil
}

func (r *Repository) GetUsers(db *gorm.DB, channelID string) ([]*models.User, error) {
	channelMembers, err := r.GetMembers(db, channelID)
	if err != nil {
		return nil, err
	}

	var users []*models.User
	for _, channelMember := range channelMembers {
		users = append(users, &channelMember.User)
	}

	return users, nil
}
