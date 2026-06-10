package service

import (
	"sespima_api/internal/repository"
	"sespima_api/models"
)

type PunishmentService struct{ repo repository.PunishmentRepository }

func NewPunishmentService(repo repository.PunishmentRepository) *PunishmentService {
	return &PunishmentService{repo: repo}
}

func (s *PunishmentService) GetAll() ([]models.PunishmentLog, error)                       { return s.repo.GetAll() }
func (s *PunishmentService) GetByUserID(id string) ([]models.PunishmentLog, error)          { return s.repo.GetByUserID(id) }
func (s *PunishmentService) GetPending() ([]models.PunishmentLog, error)                    { return s.repo.GetPending() }
func (s *PunishmentService) Create(p *models.PunishmentLog) error                           { return s.repo.Create(p) }
func (s *PunishmentService) Update(id string, p *models.PunishmentLog) error                { return s.repo.Update(id, p) }
func (s *PunishmentService) Delete(id string) error                                         { return s.repo.Delete(id) }
func (s *PunishmentService) Approve(id string, approverID int64) error                      { return s.repo.Approve(id, approverID) }
func (s *PunishmentService) Reject(id string, approverID int64, reason string) error {
	return s.repo.Reject(id, approverID, reason)
}
