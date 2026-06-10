package repository

import (
	"time"

	"sespima_api/models"

	"gorm.io/gorm"
)

type kegiatanRepo struct{ db *gorm.DB }

func NewKegiatanRepository(db *gorm.DB) KegiatanRepository { return &kegiatanRepo{db: db} }

func (r *kegiatanRepo) GetAll() ([]models.Kegiatan, error) {
	var list []models.Kegiatan
	err := r.db.Find(&list).Error
	return list, err
}

func (r *kegiatanRepo) GetRutin() ([]models.Kegiatan, error) {
	var list []models.Kegiatan
	err := r.db.Where("is_rutin = ?", true).Find(&list).Error
	return list, err
}

func (r *kegiatanRepo) GetToday() ([]models.Kegiatan, error) {
	today := time.Now().Format("2006-01-02")
	var list []models.Kegiatan
	err := r.db.Where("tanggal_mulai = ?", today).Find(&list).Error
	return list, err
}

func (r *kegiatanRepo) GetByID(id string) (*models.Kegiatan, error) {
	var k models.Kegiatan
	err := r.db.First(&k, id).Error
	return &k, err
}

func (r *kegiatanRepo) Create(k *models.Kegiatan) error {
	return r.db.Create(k).Error
}

func (r *kegiatanRepo) Update(k *models.Kegiatan) error {
	return r.db.Save(k).Error
}

func (r *kegiatanRepo) Delete(id string) error {
	return r.db.Delete(&models.Kegiatan{}, id).Error
}
