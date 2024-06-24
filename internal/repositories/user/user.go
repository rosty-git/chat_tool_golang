package userrepository

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(db *gorm.DB, user *models.User) (*models.User, error) {
	result := db.Create(user)

	slog.Info("Created user: ", "user", user)

	return user, result.Error
}

func (r *Repository) GetByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	result := db.Where("email = ?", email).First(&user)

	return &user, result.Error
}

func (r *Repository) GetById(db *gorm.DB, id string) (*models.User, error) {
	slog.Info("GetById", "id", id)

	var user models.User
	result := db.First(&user, "id = ?", id)

	return &user, result.Error
}

func (r *Repository) CreateOrUpdateStatus(db *gorm.DB, userID string, newStatus string, manual bool, dndEndTime string) (*models.Status, error) {
	slog.Info("CreateOrUpdateStatus", "user_id", userID, "newStatus", newStatus, "manual", manual, "dndEndTime", dndEndTime)

	var oldStatus models.Status
	err := db.Model(&models.Status{}).
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

		result := db.Model(&status).Clauses(clause.Returning{}).Create(statusUpdate)

		slog.Info("Created status: ", "status", status)

		return &status, result.Error
	} else {
		statusUpdateMap := map[string]interface{}{
			"user_id":          statusUpdate.UserID,
			"status":           statusUpdate.Status,
			"manual":           statusUpdate.Manual,
			"last_activity_at": statusUpdate.LastActivityAt,
		}

		if manual && newStatus == "dnd" {
			statusUpdateMap["prev_status"] = oldStatus.Status
		}

		result := db.Model(&status).Where("user_id = ?", userID).Updates(statusUpdateMap)

		slog.Info("Updated status: ", "status", status)

		return &status, result.Error
	}
}

func (r *Repository) GetStatus(db *gorm.DB, userID string) (*models.Status, error) {
	var status models.Status
	err := db.Model(&models.Status{}).
		Where("user_id = ?", userID).
		Find(&status).
		Error
	if err != nil {
		return nil, err
	}

	return &status, err
}

func (r *Repository) GetNotUpdatedStatuses(db *gorm.DB) ([]*models.Status, error) {
	cutoffTime := time.Now().Add(-5 * time.Minute)

	// Prepare the query
	var statuses []*models.Status
	err := db.Where("status IN ? AND `manual` = ? AND last_activity_at < ?", []string{"online", "away"}, false, cutoffTime).Find(&statuses).Error
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func (r *Repository) Search(db *gorm.DB, userID string, text string) ([]*models.User, error) {
	var users []*models.User
	err := db.Where("id != ?", userID).
		Where("MATCH(name) AGAINST(? IN NATURAL LANGUAGE MODE)", text).
		Find(&users).Error

	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	return users, nil
}
