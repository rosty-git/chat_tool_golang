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

func (r *Repository) IncrementTotalMsgCount(db *gorm.DB, channelID string) error {
	return db.Model(&models.Channel{}).Where("id = ?", channelID).Update("total_msg_count", gorm.Expr("total_msg_count + ?", 1)).Error
}

func (r *Repository) MarkAsRead(db *gorm.DB, channelID string, userID string) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var channel models.Channel

		if err := db.First(&channel, "id = ?", channelID).Error; err != nil {
			return err
		}

		slog.Debug("Channel repo", "channel", channel)

		if err := db.Model(&models.ChannelMembers{}).Where("user_id = ? AND channel_id = ?", userID, channelID).Update("msg_count", channel.TotalMsgCount).Error; err != nil {
			return err
		}

		return nil
	})
	return err
}

func (r *Repository) GetUnreadCount(db *gorm.DB, channelID string, userID string) (uint64, error) {
	var count uint64

	err := db.Transaction(func(tx *gorm.DB) error {
		var channel models.Channel

		if err := db.First(&channel, "id = ?", channelID).Error; err != nil {
			return err
		}

		var channelMember models.ChannelMembers

		if err := db.First(&channelMember, "user_id = ? AND channel_id = ?", userID, channelID).Error; err != nil {
			return err
		}

		count = channel.TotalMsgCount - channelMember.MsgCount

		return nil
	})
	return count, err
}
