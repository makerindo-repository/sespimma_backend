package repository

import (
	"time"

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

// GetRecapBySerdik builds the strict attendance recap across every past zone:
//   - check-in exists       → hadir (or telat when is_late)
//   - approved izin covers   → izin
//   - missed & cutoff passed → tk (Tanpa Keterangan)
//   - missed & still open    → omitted (not yet a miss)
func (r *absensiRepo) GetRecapBySerdik(serdikID uint64) ([]RecapRecord, *AbsensiSummary, error) {
	// 1. Past zones (kegiatan with a real geofence location).
	type kegRow struct {
		ID          uint64
		Nama        string
		NamaLokasi  string
		ScheduledAt time.Time
		CutoffAt    time.Time
	}
	var kegs []kegRow
	if err := r.db.Raw(`
		SELECT id, nama, nama_lokasi,
			(tanggal_mulai::date + waktu_mulai::time) AS scheduled_at,
			(tanggal_mulai::date + COALESCE(cutoff_time, batas_waktu_penugasan)::time) AS cutoff_at
		FROM kegiatan
		WHERE location IS NOT NULL AND tanggal_mulai::date <= CURRENT_DATE
		ORDER BY scheduled_at DESC`).Scan(&kegs).Error; err != nil {
		return nil, nil, err
	}

	// 2. This serdik's check-ins, keyed by zone.
	type absRow struct {
		ID         uint64
		KegiatanID uint64
		Status     string
		IsLate     bool
		Datetime   time.Time
		Method     string
	}
	var absRows []absRow
	if err := r.db.Raw(`
		SELECT id, kegiatan_id, status, is_late, datetime, method
		FROM absensi WHERE serdik_id = ?`, serdikID).Scan(&absRows).Error; err != nil {
		return nil, nil, err
	}
	byKeg := make(map[uint64]absRow, len(absRows))
	for _, a := range absRows {
		byKeg[a.KegiatanID] = a
	}

	// 3. Approved leave windows (date-range granularity).
	type izinRow struct {
		StartTime time.Time
		EndTime   time.Time
	}
	var izins []izinRow
	_ = r.db.Raw(`
		SELECT start_time, end_time FROM izin_requests
		WHERE serdik_id = ? AND status = 'disetujui'`, serdikID).Scan(&izins).Error

	coveredByIzin := func(at time.Time) bool {
		d := at.Truncate(24 * time.Hour)
		for _, iz := range izins {
			start := iz.StartTime.Truncate(24 * time.Hour)
			end := iz.EndTime.Truncate(24 * time.Hour)
			if !d.Before(start) && !d.After(end) {
				return true
			}
		}
		return false
	}

	now := time.Now()
	records := make([]RecapRecord, 0, len(kegs))
	summary := &AbsensiSummary{}

	for _, k := range kegs {
		rec := RecapRecord{
			KegiatanID: k.ID,
			Method:     "-",
			Kegiatan:   RecapKegiatan{Nama: k.Nama, NamaLokasi: k.NamaLokasi},
		}

		if a, ok := byKeg[k.ID]; ok {
			id := a.ID
			rec.ID = &id
			rec.Datetime = a.Datetime
			rec.Method = a.Method
			rec.IsLate = a.IsLate
			switch a.Status {
			case "hadir":
				if a.IsLate {
					rec.Status = "telat"
					summary.LateCount++
				} else {
					rec.Status = "hadir"
					summary.PresentCount++
				}
			case "izin", "sakit":
				rec.Status = "izin"
				summary.PermissionCount++
			default: // tk / alpha
				rec.Status = "tk"
				summary.AbsentCount++
			}
		} else if coveredByIzin(k.ScheduledAt) {
			rec.Datetime = k.ScheduledAt
			rec.Status = "izin"
			summary.PermissionCount++
		} else if now.After(k.CutoffAt) {
			rec.Datetime = k.ScheduledAt
			rec.Status = "tk"
			summary.AbsentCount++
		} else {
			// Zone still open and not yet attended — not a miss yet.
			continue
		}

		records = append(records, rec)
	}

	return records, summary, nil
}

func (r *absensiRepo) GetSummaryBySerdik(serdikID uint64) (*AbsensiSummary, error) {
	type row struct {
		Status string
		IsLate bool
		Count  int
	}
	var rows []row
	err := r.db.Raw(`
		SELECT status, is_late, COUNT(*) AS count FROM absensi
		WHERE serdik_id = ? GROUP BY status, is_late`, serdikID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	s := &AbsensiSummary{}
	for _, r := range rows {
		switch r.Status {
		case "hadir":
			if r.IsLate {
				s.LateCount += r.Count
			} else {
				s.PresentCount += r.Count
			}
		case "izin", "sakit":
			s.PermissionCount += r.Count
		case "tk", "alpha":
			s.AbsentCount += r.Count
		}
	}
	return s, nil
}
