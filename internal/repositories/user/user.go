package userrepository

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(user *models.User) (*models.User, error) {
	result := r.db.Create(user)

	slog.Info("Created user: ", user)

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

func (r *Repository) CreateOrUpdateStatus(userID string, newStatus string, manual bool, dndEndTime time.Time) (*models.Status, error) {
	slog.Info("CreateOrUpdateStatus", "user_id", userID, "newStatus", newStatus, "manual", manual, "dndEndTime", dndEndTime)

	var oldStatus *models.Status
	err := r.db.Model(&models.Status{}).
		Where("user_id = ?", userID).
		Find(&oldStatus).
		Error
	if err != nil {
		return nil, err
	}

	var status *models.Status

	statusUpdate := &models.Status{
		UserID:         userID,
		Status:         newStatus,
		Manual:         manual,
		LastActivityAt: time.Now(),
		DNDEndTime:     dndEndTime,
		PrevStatus:     "",
	}

	if oldStatus.UserID == "" {
		result := r.db.Model(status).Create(statusUpdate)

		return status, result.Error
	} else {
		statusUpdate.PrevStatus = oldStatus.Status

		result := r.db.Model(status).Where("user_id = ?", userID).Updates(statusUpdate)

		return status, result.Error
	}
}
