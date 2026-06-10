package repository

import (
	"sespima_api/models"

	"gorm.io/gorm"
)

type locationLogRepo struct{ db *gorm.DB }

func NewLocationLogRepository(db *gorm.DB) LocationLogRepository {
	return &locationLogRepo{db: db}
}

func (r *locationLogRepo) Create(log *models.UserLocationLog) error {
	return r.db.Create(log).Error
}

func (r *locationLogRepo) GetLatestByUser(userID int64) (*models.UserLocationLog, error) {
	var log models.UserLocationLog
	err := r.db.Where("user_id = ?", userID).
		Order("recorded_at DESC").
		First(&log).Error
	return &log, err
}
