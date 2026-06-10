package repository

import (
	"sespima_api/models"

	"gorm.io/gorm"
)

type userRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepo{db: db} }

func (r *userRepo) FindByNrpOrNip(nrpNip string) (*models.User, error) {
	var u models.User
	err := r.db.Where("nrp = ? OR nip = ?", nrpNip, nrpNip).First(&u).Error
	return &u, err
}

func (r *userRepo) FindByID(id uint) (*models.User, error) {
	var u models.User
	err := r.db.First(&u, id).Error
	return &u, err
}

func (r *userRepo) Save(user *models.User) error {
	return r.db.Omit("location", "Location").Save(user).Error
}

func (r *userRepo) UpdateLocation(userID interface{}, lat, lng float64) error {
	return r.db.Exec(
		`UPDATE users SET location = ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography WHERE id = ?`,
		lng, lat, userID,
	).Error
}

func (r *userRepo) GetAllWithProfiles() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Pimpinan").Preload("Patun").Preload("Korsis").
		Preload("Serdik").Preload("Gadik").Find(&users).Error
	return users, err
}

func (r *userRepo) GetUserLocation(id string) (*UserLocationResult, error) {
	var result UserLocationResult
	err := r.db.Raw(`
		SELECT id, email,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude
		FROM users WHERE id = ?`, id).Scan(&result).Error
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, err
}

func (r *userRepo) GetAllUserLocations() ([]UserLocationResult, error) {
	var results []UserLocationResult
	err := r.db.Raw(`
		SELECT id, email, role,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude
		FROM users WHERE location IS NOT NULL ORDER BY id`).Scan(&results).Error
	return results, err
}
