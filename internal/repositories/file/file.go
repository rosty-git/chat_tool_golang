package filerepository

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetTmpById(db *gorm.DB, id string) (*models.TmpFile, error) {
	var file models.TmpFile

	err := db.First(&file, "id = ?", id).Error

	return &file, err
}

func (r *Repository) GetById(db *gorm.DB, id string) (*models.File, error) {
	var file models.File

	err := db.First(&file, "id = ?", id).Error

	return &file, err
}

func (r *Repository) CreateTmp(db *gorm.DB, fileName, fileType string, size uint64) (*models.TmpFile, error) {
	file := &models.TmpFile{
		Name: fileName,
		Type: fileType,
		Size: size,
	}
	result := db.Create(file)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.GetTmpById(db, file.ID)
}

func (r *Repository) Create(db *gorm.DB, id, name, fileType, postID, s3Key string, size uint64) (*models.File, error) {
	file := &models.File{
		BaseModel: models.BaseModel{
			ID: id,
		},
		Name:   name,
		Type:   fileType,
		Size:   size,
		PostID: postID,
		S3Key:  s3Key,
	}
	result := db.Create(file)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.GetById(db, file.ID)
}

func (r *Repository) UpdateTmp(db *gorm.DB, id string, m map[string]interface{}) error {
	err := db.Model(&models.TmpFile{BaseModel: models.BaseModel{ID: id}}).Updates(m).Error

	return err
}

func (r *Repository) DeleteTmp(db *gorm.DB, id string) error {
	err := db.Where("id = ?", id).Delete(&models.TmpFile{}).Error

	return err
}
