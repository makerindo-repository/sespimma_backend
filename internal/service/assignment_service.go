package service

import (
	"sespima_api/internal/repository"
	"sespima_api/models"
)

type AssignmentService struct{ repo repository.AssignmentRepository }

func NewAssignmentService(repo repository.AssignmentRepository) *AssignmentService {
	return &AssignmentService{repo: repo}
}

func (s *AssignmentService) GetAll(pokjarID *int) ([]models.Assignment, error) {
	return s.repo.GetAll(pokjarID)
}

func (s *AssignmentService) GetByID(id int64) (*models.Assignment, error) {
	return s.repo.GetByID(id)
}

func (s *AssignmentService) Create(a *models.Assignment) error { return s.repo.Create(a) }

func (s *AssignmentService) Update(id int64, input *models.Assignment) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	existing.Judul = input.Judul
	existing.Deskripsi = input.Deskripsi
	existing.JenisTugas = input.JenisTugas
	existing.TurunanTugas = input.TurunanTugas
	existing.Mapel = input.Mapel
	existing.Deadline = input.Deadline
	existing.TargetPokjarID = input.TargetPokjarID
	existing.Instruksi = input.Instruksi
	existing.Status = input.Status
	return s.repo.Update(existing)
}

func (s *AssignmentService) Delete(id int64) error { return s.repo.Delete(id) }

func (s *AssignmentService) GetSubmissions(assignmentID int64) ([]models.Submission, error) {
	return s.repo.GetSubmissions(assignmentID)
}

func (s *AssignmentService) Submit(sub *models.Submission) error {
	return s.repo.CreateSubmission(sub)
}

func (s *AssignmentService) GradeSubmission(sub *models.Submission) error {
	return s.repo.UpdateSubmission(sub)
}
