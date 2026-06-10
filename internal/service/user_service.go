package service

import (
	"fmt"

	"sespima_api/internal/repository"
	"sespima_api/models"
)

type UserService struct {
	users       repository.UserRepository
	locationLog repository.LocationLogRepository
}

func NewUserService(users repository.UserRepository, locationLog repository.LocationLogRepository) *UserService {
	return &UserService{users: users, locationLog: locationLog}
}

func (s *UserService) GetAllWithProfiles() ([]models.User, error) {
	return s.users.GetAllWithProfiles()
}

func (s *UserService) GetUserLocation(id string) (*repository.UserLocationResult, error) {
	result, err := s.users.GetUserLocation(id)
	if err != nil {
		return nil, fmt.Errorf("user tidak ditemukan")
	}
	return result, nil
}

func (s *UserService) GetAllUserLocations() ([]repository.UserLocationResult, error) {
	return s.users.GetAllUserLocations()
}

func (s *UserService) UpdateLocation(userID interface{}, lat, lng float64) error {
	return s.users.UpdateLocation(userID, lat, lng)
}

func (s *UserService) LogLocation(log *models.UserLocationLog) error {
	return s.locationLog.Create(log)
}

func (s *UserService) GetByID(userID uint) (*models.User, error) {
	return s.users.FindByID(userID)
}

func (s *UserService) UpdateProfilePhoto(userID uint, photoURL string) error {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan")
	}
	user.ProfilePhoto = &photoURL
	return s.users.Save(user)
}

func (s *UserService) ClearProfilePhoto(userID uint) error {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan")
	}
	user.ProfilePhoto = nil
	return s.users.Save(user)
}
