package repository

import (
	"sespima_api/models"

	"gorm.io/gorm"
)

type rewardRepo struct{ db *gorm.DB }

func NewRewardRepository(db *gorm.DB) RewardRepository { return &rewardRepo{db: db} }

func (r *rewardRepo) GetAll() ([]models.UserReward, error) {
	var list []models.UserReward
	err := r.db.Preload("RewardItem").Find(&list).Error
	return list, err
}

func (r *rewardRepo) GetPending() ([]models.UserReward, error) {
	var list []models.UserReward
	err := r.db.Where("approved_by IS NULL AND status = 'pending'").
		Preload("RewardItem").Find(&list).Error
	return list, err
}

func (r *rewardRepo) Create(rw *models.UserReward) error {
	// Maker-Checker: a reward assigned by patun/gadik is always created as
	// pending and unapproved, regardless of what the client sent.
	rw.Status = "pending"
	rw.ApprovedBy = nil
	return r.db.Create(rw).Error
}

func (r *rewardRepo) Update(id string, input *models.UserReward) error {
	return r.db.Model(&models.UserReward{}).Where("id = ?", id).Updates(input).Error
}

func (r *rewardRepo) Delete(id string) error {
	return r.db.Delete(&models.UserReward{}, id).Error
}

func (r *rewardRepo) Approve(id string, approverID int64) error {
	return r.db.Model(&models.UserReward{}).Where("id = ? AND status = 'pending'", id).
		Updates(map[string]interface{}{"approved_by": approverID, "status": "approved"}).Error
}

func (r *rewardRepo) Reject(id string, approverID int64, reason string) error {
	return r.db.Model(&models.UserReward{}).Where("id = ? AND status = 'pending'", id).
		Updates(map[string]interface{}{
			"approved_by": approverID,
			"status":      "rejected",
			"notes":       reason,
		}).Error
}
