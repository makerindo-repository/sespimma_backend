package repository

import (
	"sespima_api/models"

	"gorm.io/gorm"
)

type punishmentRepo struct{ db *gorm.DB }

func NewPunishmentRepository(db *gorm.DB) PunishmentRepository { return &punishmentRepo{db: db} }

func (r *punishmentRepo) GetAll() ([]models.PunishmentLog, error) {
	var list []models.PunishmentLog
	err := r.db.Preload("PunishmentItem").Find(&list).Error
	return list, err
}

func (r *punishmentRepo) GetByUserID(userID string) ([]models.PunishmentLog, error) {
	var list []models.PunishmentLog
	err := r.db.Where("user_id = ?", userID).Preload("PunishmentItem").Find(&list).Error
	return list, err
}

func (r *punishmentRepo) Create(p *models.PunishmentLog) error {
	return r.db.Create(p).Error
}

func (r *punishmentRepo) Update(id string, input *models.PunishmentLog) error {
	return r.db.Model(&models.PunishmentLog{}).Where("id = ?", id).Updates(input).Error
}

func (r *punishmentRepo) Delete(id string) error {
	return r.db.Delete(&models.PunishmentLog{}, id).Error
}
