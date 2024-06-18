package postusecase

import (
	"log/slog"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type postService interface {
	GetByChannelId(db *gorm.DB, channelID string, limit int, before string, after string) ([]*models.Post, error)
	Create(db *gorm.DB, userID string, channelID string, message string) (*models.Post, error)
	NotifyReceivers(usersIDs []string, message interface{})
	NotifySender(userID string, message interface{})
	GetById(db *gorm.DB, id string) (*models.Post, error)
}

type channelService interface {
	GetUsers(db *gorm.DB, channelID string) ([]*models.User, error)
	IncrementTotalMsgCount(db *gorm.DB, channelID string) error
	MarkAsRead(db *gorm.DB, channelID string, userID string) error
}

type fileService interface {
	GetTmp(db *gorm.DB, id string) (*models.TmpFile, error)
	DeleteTmp(db *gorm.DB, id string) error
	Create(db *gorm.DB, id, name, fileType, postID, s3Key string, size uint64) (*models.File, error)
}

type UseCase struct {
	postService    postService
	channelService channelService
	fileService    fileService
	db             *gorm.DB
}

func NewUseCase(db *gorm.DB, postService postService, channelService channelService, fileService fileService) *UseCase {
	return &UseCase{
		db:             db,
		postService:    postService,
		channelService: channelService,
		fileService:    fileService,
	}
}

func (uc *UseCase) GetByChannelId(channelID string, limit int, before string, after string) ([]*models.Post, error) {
	return uc.postService.GetByChannelId(uc.db, channelID, limit, before, after)
}

func (uc *UseCase) Create(userID string, channelID string, message string, frontId string, filesIDs []string) (*models.Post, error) {
	var createdPost *models.Post

	err := uc.db.Transaction(func(tx *gorm.DB) error {
		cp, err := uc.postService.Create(tx, userID, channelID, message)
		if err != nil {
			slog.Error("CreatePost", "err", err)

			return nil
		}

		for _, fileID := range filesIDs {
			tmpFile, err := uc.fileService.GetTmp(tx, fileID)
			if err != nil {
				return err
			}

			_, err = uc.fileService.Create(tx, tmpFile.ID, tmpFile.Name, tmpFile.Type, cp.ID, tmpFile.S3Key, tmpFile.Size)
			if err != nil {
				return err
			}

			err = uc.fileService.DeleteTmp(tx, tmpFile.ID)
		}

		createdPost, err = uc.postService.GetById(tx, cp.ID)
		if err != nil {
			return err
		}

		if err = uc.channelService.IncrementTotalMsgCount(tx, channelID); err != nil {
			return err
		}

		if err = uc.channelService.MarkAsRead(tx, channelID, userID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	channelUsers, err := uc.channelService.GetUsers(uc.db, channelID)
	if err != nil {
		return nil, err
	}

	var usersIDs []string

	for _, channelUser := range channelUsers {
		if channelUser.ID != userID {
			usersIDs = append(usersIDs, channelUser.ID)
		}
	}

	uc.postService.NotifyReceivers(usersIDs, createdPost)

	createdPostMap := map[string]interface{}{
		"createdPost": createdPost,
		"frontId":     frontId,
	}

	uc.postService.NotifySender(userID, createdPostMap)

	return createdPost, nil
}
