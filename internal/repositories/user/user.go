package userrepository

import (
	"log/slog"
	"time"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormDb interface {
	Create(value interface{}) (tx *gorm.DB)
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Model(value interface{}) (tx *gorm.DB)
}

type Repository struct {
	db GormDb
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(user *models.User) (*models.User, error) {
	result := r.db.Create(user)

	slog.Info("Created user: ", "user", user)

	return user, result.Error
}

func (r *Repository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)

	return &user, result.Error
}

func (r *Repository) GetById(id string) (*models.User, error) {
	slog.Info("GetById", "id", id)

	var user models.User
	result := r.db.First(&user, "id = ?", id)

	return &user, result.Error
}

func (r *Repository) CreateOrUpdateStatus(userID string, newStatus string, manual bool, dndEndTime string) (*models.Status, error) {
	slog.Info("CreateOrUpdateStatus", "user_id", userID, "newStatus", newStatus, "manual", manual, "dndEndTime", dndEndTime)

	var oldStatus models.Status
	err := r.db.Model(&models.Status{}).
		Where("user_id = ?", userID).
		Find(&oldStatus).
		Error
	if err != nil {
		return nil, err
	}

	var status models.Status

	statusUpdate := &models.Status{
		UserID:         userID,
		Status:         newStatus,
		Manual:         manual,
		LastActivityAt: time.Now(),
	}

	if oldStatus.UserID == "" {
		if dndEndTime == "" {
			statusUpdate.DNDEndTime = time.Now()
		}

		result := r.db.Model(&status).Clauses(clause.Returning{}).Create(statusUpdate)

		slog.Info("Created status: ", "status", status)

		return &status, result.Error
	} else {
		if manual && newStatus == "dnd" {
			statusUpdate.PrevStatus = oldStatus.Status
		}

		result := r.db.Model(&status).Clauses(clause.Returning{}).Where("user_id = ?", userID).Updates(statusUpdate)

		slog.Info("Updated status: ", "status", status)

		return &status, result.Error
	}
}

func (r *Repository) GetStatus(userID string) (*models.Status, error) {
	var status models.Status
	err := r.db.Model(&models.Status{}).
		Where("user_id = ?", userID).
		Find(&status).
		Error
	if err != nil {
		return nil, err
	}

	return &status, err
}
