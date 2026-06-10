package service

import (
	"sespima_api/internal/repository"
	"sespima_api/models"
)

type RewardService struct{ repo repository.RewardRepository }

func NewRewardService(repo repository.RewardRepository) *RewardService {
	return &RewardService{repo: repo}
}

func (s *RewardService) GetAll() ([]models.UserReward, error)     { return s.repo.GetAll() }
func (s *RewardService) GetPending() ([]models.UserReward, error)  { return s.repo.GetPending() }
func (s *RewardService) Create(r *models.UserReward) error          { return s.repo.Create(r) }
func (s *RewardService) Update(id string, r *models.UserReward) error { return s.repo.Update(id, r) }
func (s *RewardService) Delete(id string) error                     { return s.repo.Delete(id) }
func (s *RewardService) Approve(id string, approverID int64) error  { return s.repo.Approve(id, approverID) }
