package repository

import (
	"sespima_api/models"

	"gorm.io/gorm"
)

type sociometryRepo struct{ db *gorm.DB }

func NewSociometryRepository(db *gorm.DB) SociometryRepository { return &sociometryRepo{db: db} }

func (r *sociometryRepo) GetPeriods(pokjarID *int) ([]models.SociometryPeriod, error) {
	var list []models.SociometryPeriod
	q := r.db.Preload("Pokjar")
	if pokjarID != nil {
		q = q.Where("pokjar_id = ?", *pokjarID)
	}
	err := q.Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *sociometryRepo) CreatePeriod(p *models.SociometryPeriod) error {
	return r.db.Create(p).Error
}

func (r *sociometryRepo) ClosePeriod(id int64) error {
	return r.db.Model(&models.SociometryPeriod{}).
		Where("id = ?", id).Update("is_active", false).Error
}

func (r *sociometryRepo) CreateEvaluation(e *models.SociometryEvaluation) error {
	return r.db.Create(e).Error
}

func (r *sociometryRepo) GetEvaluationsByPeriod(periodID int64) ([]models.SociometryEvaluation, error) {
	var list []models.SociometryEvaluation
	err := r.db.Where("period_id = ?", periodID).
		Preload("Evaluator").Preload("Evaluated").Find(&list).Error
	return list, err
}

func (r *sociometryRepo) GetAverageScore(serdikID int64, periodType string) (float64, error) {
	var avg float64
	err := r.db.Raw(`
		SELECT COALESCE(AVG(e.score), 0)
		FROM sociometry_evaluations e
		JOIN sociometry_periods p ON p.id = e.period_id
		WHERE e.evaluated_serdik_id = ? AND p.period_type = ?`,
		serdikID, periodType).Scan(&avg).Error
	return avg, err
}
