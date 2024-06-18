package fileservice

import (
	"log/slog"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type fileRepository interface {
	CreateTmp(db *gorm.DB, fileName, fileType string, size uint64) (*models.TmpFile, error)
	Create(db *gorm.DB, id, name, fileType, postID, s3Key string, size uint64) (*models.File, error)
	UpdateTmp(db *gorm.DB, id string, m map[string]interface{}) error
	DeleteTmp(db *gorm.DB, id string) error
	GetTmpById(db *gorm.DB, id string) (*models.TmpFile, error)
}

type Service struct {
	fileRepository fileRepository
}

func NewService(fileRepository fileRepository) *Service {
	return &Service{
		fileRepository: fileRepository,
	}
}

func (fs *Service) CreateTmp(db *gorm.DB, name, fileType string, size uint64) (*models.TmpFile, error) {
	return fs.fileRepository.CreateTmp(db, name, fileType, size)
}

func (fs *Service) Create(db *gorm.DB, id, name, fileType, postID, s3Key string, size uint64) (*models.File, error) {
	return fs.fileRepository.Create(db, id, name, fileType, postID, s3Key, size)
}

func (fs *Service) UpdateTmp(db *gorm.DB, id string, m map[string]interface{}) error {
	return fs.fileRepository.UpdateTmp(db, id, m)
}

func (fs *Service) SetS3Key(db *gorm.DB, id, s3Key string) (*models.TmpFile, error) {
	slog.Info("FileService SetS3Key", "id", id, "s3Key", s3Key)

	file, err := fs.fileRepository.GetTmpById(db, id)
	if err != nil {
		return nil, err
	}

	if file.S3Key == "" {
		err := fs.fileRepository.UpdateTmp(db, id, map[string]interface{}{"s3_key": s3Key})
		if err != nil {
			return nil, err
		}

		return fs.fileRepository.GetTmpById(db, id)
	}

	return file, err
}

func (fs *Service) DeleteTmp(db *gorm.DB, id string) error {
	return fs.fileRepository.DeleteTmp(db, id)
}

func (fs *Service) GetTmp(db *gorm.DB, id string) (*models.TmpFile, error) {
	return fs.fileRepository.GetTmpById(db, id)
}
