package repository

import (
	"sespima_api/models"

	"gorm.io/gorm"
)

type absensiRepo struct{ db *gorm.DB }

func NewAbsensiRepository(db *gorm.DB) AbsensiRepository { return &absensiRepo{db: db} }

func (r *absensiRepo) Create(a *models.Absensi) error {
	return r.db.Create(a).Error
}

func (r *absensiRepo) GetBySerdikID(serdikID uint64) ([]models.Absensi, error) {
	var list []models.Absensi
	err := r.db.Where("serdik_id = ?", serdikID).
		Preload("Kegiatan").Order("datetime DESC").Find(&list).Error
	return list, err
}

func (r *absensiRepo) GetByUserID(userID uint64) ([]models.Absensi, error) {
	var list []models.Absensi
	err := r.db.Where("user_id = ?", userID).
		Order("datetime DESC").Find(&list).Error
	return list, err
}

func (r *absensiRepo) GetByKegiatanID(kegiatanID uint64) ([]models.Absensi, error) {
	var list []models.Absensi
	err := r.db.Where("kegiatan_id = ?", kegiatanID).
		Preload("Serdik").Find(&list).Error
	return list, err
}

func (r *absensiRepo) GetSummaryBySerdik(serdikID uint64) (*AbsensiSummary, error) {
	type row struct {
		Status string
		Count  int
	}
	var rows []row
	err := r.db.Raw(`
		SELECT status, COUNT(*) AS count FROM absensi
		WHERE serdik_id = ? GROUP BY status`, serdikID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	s := &AbsensiSummary{}
	for _, r := range rows {
		switch r.Status {
		case "hadir":
			s.PresentCount = r.Count
		case "izin", "sakit":
			s.PermissionCount += r.Count
		case "tk", "alpha":
			s.AbsentCount += r.Count
		}
	}
	return s, nil
}
