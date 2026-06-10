package repository

import (
	"sespima_api/models"

	"gorm.io/gorm"
)

// ---- Akademik ----

type akademikRepo struct{ db *gorm.DB }

func NewAkademikRepository(db *gorm.DB) AkademikRepository { return &akademikRepo{db: db} }

func (r *akademikRepo) GetAllComponents() ([]models.Akademik, error) {
	var list []models.Akademik
	err := r.db.Preload("Children.Children.Children.Children").
		Where("parent_id IS NULL").Find(&list).Error
	return list, err
}

func (r *akademikRepo) GetPenilaianBySerdik(serdikID string) ([]models.PenilaianAkademik, error) {
	var list []models.PenilaianAkademik
	err := r.db.Where("serdik_id = ?", serdikID).Find(&list).Error
	return list, err
}

func (r *akademikRepo) CreatePenilaian(p *models.PenilaianAkademik) error {
	return r.db.Create(p).Error
}

func (r *akademikRepo) UpdatePenilaian(id string, nilai float64, catatan *string) error {
	return r.db.Model(&models.PenilaianAkademik{}).Where("id = ?", id).
		Updates(map[string]interface{}{"nilai": nilai, "catatan": catatan}).Error
}

func (r *akademikRepo) DeletePenilaian(id string) error {
	return r.db.Delete(&models.PenilaianAkademik{}, id).Error
}

// ---- Jasmani ----

type jasmaniRepo struct{ db *gorm.DB }

func NewJasmaniRepository(db *gorm.DB) JasmaniRepository { return &jasmaniRepo{db: db} }

func (r *jasmaniRepo) GetAllComponents() ([]models.Jasmani, error) {
	var list []models.Jasmani
	err := r.db.Preload("Children.Children.Children").
		Where("parent_id IS NULL").Find(&list).Error
	return list, err
}

func (r *jasmaniRepo) GetPenilaianBySerdik(serdikID string) ([]models.PenilaianJasmani, error) {
	var list []models.PenilaianJasmani
	err := r.db.Where("serdik_id = ?", serdikID).Find(&list).Error
	return list, err
}

func (r *jasmaniRepo) CreatePenilaian(p *models.PenilaianJasmani) error {
	return r.db.Create(p).Error
}

func (r *jasmaniRepo) UpdatePenilaian(id string, nilai float64, catatan *string) error {
	return r.db.Model(&models.PenilaianJasmani{}).Where("id = ?", id).
		Updates(map[string]interface{}{"nilai": nilai, "catatan": catatan}).Error
}

func (r *jasmaniRepo) DeletePenilaian(id string) error {
	return r.db.Delete(&models.PenilaianJasmani{}, id).Error
}

// ---- Mental ----

type mentalRepo struct{ db *gorm.DB }

func NewMentalRepository(db *gorm.DB) MentalRepository { return &mentalRepo{db: db} }

func (r *mentalRepo) GetAllComponents() ([]models.Mental, error) {
	var list []models.Mental
	err := r.db.Preload("Children.Children.Children").
		Where("parent_id IS NULL").Find(&list).Error
	return list, err
}

func (r *mentalRepo) GetPenilaianBySerdik(serdikID string) ([]models.PenilaianMental, error) {
	var list []models.PenilaianMental
	err := r.db.Where("serdik_id = ?", serdikID).Find(&list).Error
	return list, err
}

func (r *mentalRepo) CreatePenilaian(p *models.PenilaianMental) error {
	return r.db.Create(p).Error
}

func (r *mentalRepo) UpdatePenilaian(id string, nilai float64, catatan *string) error {
	return r.db.Model(&models.PenilaianMental{}).Where("id = ?", id).
		Updates(map[string]interface{}{"nilai": nilai, "catatan": catatan}).Error
}

func (r *mentalRepo) DeletePenilaian(id string) error {
	return r.db.Delete(&models.PenilaianMental{}, id).Error
}
