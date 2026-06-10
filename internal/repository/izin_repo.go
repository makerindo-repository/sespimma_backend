package repository

import (
	"time"

	"sespima_api/models"

	"gorm.io/gorm"
)

type izinRepo struct{ db *gorm.DB }

func NewIzinRepository(db *gorm.DB) IzinRepository { return &izinRepo{db: db} }

func (r *izinRepo) Create(req *models.IzinRequest) error {
	return r.db.Create(req).Error
}

func (r *izinRepo) GetBySerdikID(serdikID int64) ([]models.IzinRequest, error) {
	var list []models.IzinRequest
	err := r.db.Where("serdik_id = ?", serdikID).
		Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *izinRepo) GetPending() ([]models.IzinRequest, error) {
	var list []models.IzinRequest
	err := r.db.Where("status = 'pending'").
		Preload("Serdik").Order("created_at ASC").Find(&list).Error
	return list, err
}

func (r *izinRepo) UpdateStatus(id int64, status string, reviewerID int64, reason *string) error {
	now := time.Now()
	return r.db.Model(&models.IzinRequest{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           status,
			"reviewed_by":      reviewerID,
			"reviewed_at":      now,
			"rejection_reason": reason,
		}).Error
}
