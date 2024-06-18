package fileusecase

import (
	"log/slog"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type fileService interface {
	CreateTmp(db *gorm.DB, name, fileType string, size uint64) (*models.TmpFile, error)
	DeleteTmp(db *gorm.DB, id string) error
	SetS3Key(db *gorm.DB, fileId, s3Key string) (*models.TmpFile, error)
	GetTmp(db *gorm.DB, id string) (*models.TmpFile, error)
}

type s3Storage interface {
	GetPresignUrl(key string) (string, error)
	Delete(key string) error
}

type UseCase struct {
	db          *gorm.DB
	fileService fileService
	s3Storage   s3Storage
}

func NewUseCase(db *gorm.DB, fileService fileService, storage s3Storage) *UseCase {
	return &UseCase{
		db:          db,
		fileService: fileService,
		s3Storage:   storage,
	}
}

func (fu *UseCase) CreateTmp(name, fileType string, size uint64) (*models.TmpFile, error) {
	return fu.fileService.CreateTmp(fu.db, name, fileType, size)
}

func (fu *UseCase) SetS3Key(id, s3Key string) (*models.TmpFile, error) {
	slog.Info("FileUseCase SetS3Key", "id", id, "s3Key", s3Key)

	return fu.fileService.SetS3Key(fu.db, id, s3Key)
}

func (fu *UseCase) DeleteTmp(id string) error {
	file, err := fu.fileService.GetTmp(fu.db, id)
	if err != nil {
		return err
	}

	err = fu.s3Storage.Delete(file.S3Key)
	if err != nil {
		slog.Error("s3Storage.Delete", "err", err)
	}

	return fu.fileService.DeleteTmp(fu.db, id)
}

func (fu *UseCase) GetPresignUrl(key string) (string, error) {
	return fu.s3Storage.GetPresignUrl(key)
}
