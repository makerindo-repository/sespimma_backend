package repository

import (
	"sespima_api/models"

	"gorm.io/gorm"
)

type healthRepo struct{ db *gorm.DB }

func NewHealthRepository(db *gorm.DB) HealthRepository { return &healthRepo{db: db} }

func (r *healthRepo) GetBySerdikID(serdikID int64) (*models.SerdikHealthData, error) {
	var d models.SerdikHealthData
	err := r.db.Where("serdik_id = ?", serdikID).
		Preload("Records").First(&d).Error
	return &d, err
}

func (r *healthRepo) UpsertHealthData(d *models.SerdikHealthData) error {
	return r.db.Save(d).Error
}

func (r *healthRepo) CreateRecord(rec *models.HealthRecord) error {
	return r.db.Create(rec).Error
}

func (r *healthRepo) GetRecords(healthDataID int64) ([]models.HealthRecord, error) {
	var list []models.HealthRecord
	err := r.db.Where("serdik_health_data_id = ?", healthDataID).
		Order("recorded_at DESC").Find(&list).Error
	return list, err
}
