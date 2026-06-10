package repository

import (
	"sespima_api/models"

	"gorm.io/gorm"
)

type assignmentRepo struct{ db *gorm.DB }

func NewAssignmentRepository(db *gorm.DB) AssignmentRepository { return &assignmentRepo{db: db} }

func (r *assignmentRepo) GetAll(pokjarID *int) ([]models.Assignment, error) {
	var list []models.Assignment
	q := r.db.Preload("Creator").Preload("Pokjar")
	if pokjarID != nil {
		q = q.Where("target_pokjar_id = ? OR target_pokjar_id IS NULL", *pokjarID)
	}
	err := q.Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *assignmentRepo) GetByID(id int64) (*models.Assignment, error) {
	var a models.Assignment
	err := r.db.Preload("Creator").Preload("Pokjar").First(&a, id).Error
	return &a, err
}

func (r *assignmentRepo) Create(a *models.Assignment) error {
	return r.db.Create(a).Error
}

func (r *assignmentRepo) Update(a *models.Assignment) error {
	return r.db.Save(a).Error
}

func (r *assignmentRepo) Delete(id int64) error {
	return r.db.Delete(&models.Assignment{}, id).Error
}

func (r *assignmentRepo) GetSubmissions(assignmentID int64) ([]models.Submission, error) {
	var list []models.Submission
	err := r.db.Where("assignment_id = ?", assignmentID).
		Preload("Serdik").Find(&list).Error
	return list, err
}

func (r *assignmentRepo) CreateSubmission(s *models.Submission) error {
	return r.db.Create(s).Error
}

func (r *assignmentRepo) UpdateSubmission(s *models.Submission) error {
	return r.db.Save(s).Error
}
