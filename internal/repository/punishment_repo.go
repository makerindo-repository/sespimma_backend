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

func (r *punishmentRepo) GetPending() ([]models.PunishmentLog, error) {
	var list []models.PunishmentLog
	err := r.db.Where("status = 'pending'").
		Preload("PunishmentItem").Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *punishmentRepo) Create(p *models.PunishmentLog) error {
	// Maker-Checker: a punishment logged by patun/gadik is always created as
	// pending and unapproved, regardless of what the client sent.
	p.Status = "pending"
	p.ApprovedBy = nil
	return r.db.Create(p).Error
}

func (r *punishmentRepo) Update(id string, input *models.PunishmentLog) error {
	return r.db.Model(&models.PunishmentLog{}).Where("id = ?", id).Updates(input).Error
}

func (r *punishmentRepo) Delete(id string) error {
	return r.db.Delete(&models.PunishmentLog{}, id).Error
}

func (r *punishmentRepo) Approve(id string, approverID int64) error {
	return r.db.Model(&models.PunishmentLog{}).Where("id = ? AND status = 'pending'", id).
		Updates(map[string]interface{}{
			"approved_by": approverID,
			"status":      "approved",
			"reviewed_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *punishmentRepo) Reject(id string, approverID int64, reason string) error {
	return r.db.Model(&models.PunishmentLog{}).Where("id = ? AND status = 'pending'", id).
		Updates(map[string]interface{}{
			"approved_by":      approverID,
			"status":           "rejected",
			"reviewed_at":      gorm.Expr("NOW()"),
			"rejection_reason": reason,
		}).Error
}
