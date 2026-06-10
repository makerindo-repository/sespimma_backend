package service

import (
	"errors"
	"fmt"
	"time"

	"sespima_api/internal/repository"
	"sespima_api/models"
)

type KegiatanService struct {
	repo    repository.KegiatanRepository
	absensi repository.AbsensiRepository
	izin    repository.IzinRepository
}

func NewKegiatanService(
	repo repository.KegiatanRepository,
	absensi repository.AbsensiRepository,
	izin repository.IzinRepository,
) *KegiatanService {
	return &KegiatanService{repo: repo, absensi: absensi, izin: izin}
}

func (s *KegiatanService) GetAll() ([]models.Kegiatan, error)  { return s.repo.GetAll() }
func (s *KegiatanService) GetRutin() ([]models.Kegiatan, error) { return s.repo.GetRutin() }
func (s *KegiatanService) GetToday() ([]models.Kegiatan, error) { return s.repo.GetToday() }

func (s *KegiatanService) Create(k *models.Kegiatan) error { return s.repo.Create(k) }

func (s *KegiatanService) Update(id string, input *models.Kegiatan) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	existing.Nama = input.Nama
	existing.NamaLokasi = input.NamaLokasi
	existing.TanggalMulai = input.TanggalMulai
	existing.WaktuMulai = input.WaktuMulai
	existing.WaktuSelesai = input.WaktuSelesai
	existing.BatasWaktuPenugasan = input.BatasWaktuPenugasan
	existing.IsRutin = input.IsRutin
	existing.IsPelatihan = input.IsPelatihan
	existing.QRCode = input.QRCode
	existing.Radius = input.Radius
	existing.Location = input.Location
	return s.repo.Update(existing)
}

func (s *KegiatanService) Delete(id string) error { return s.repo.Delete(id) }

// ---- Absensi ----

func (s *KegiatanService) CheckIn(a *models.Absensi) error {
	if a.KegiatanID == nil {
		return errors.New("kegiatan_id is required")
	}
	if a.Latitude == nil || a.Longitude == nil {
		return errors.New("latitude and longitude are required")
	}

	// 1. Fetch Kegiatan (Zone)
	kegiatanIDStr := fmt.Sprintf("%d", *a.KegiatanID)
	k, err := s.repo.GetByID(kegiatanIDStr)
	if err != nil {
		return errors.New("kegiatan tidak ditemukan")
	}

	// 2. Validate Distance
	dist, err := s.repo.GetDistance(*a.KegiatanID, *a.Latitude, *a.Longitude)
	if err != nil {
		return errors.New("gagal menghitung jarak lokasi")
	}

	if k.Radius != nil && dist > *k.Radius {
		return errors.New("lokasi berada di luar radius zona absensi")
	}

	// 3. Validate Time
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	// Build deadline and cutoff times using TanggalMulai + time strings
	dateStr := k.TanggalMulai.Format("2006-01-02")
	deadlineStr := fmt.Sprintf("%s %s", dateStr, k.BatasWaktuPenugasan)
	deadline, _ := time.ParseInLocation("2006-01-02 15:04:05", deadlineStr, loc)

	cutoffStr := fmt.Sprintf("%s 23:59:59", dateStr)
	if k.CutoffTime != nil {
		cutoffStr = fmt.Sprintf("%s %s", dateStr, *k.CutoffTime)
	}
	cutoff, _ := time.ParseInLocation("2006-01-02 15:04:05", cutoffStr, loc)

	if now.After(cutoff) {
		a.Status = models.AttendanceStatusTK
		a.IsLate = true
	} else if now.After(deadline) {
		a.Status = models.AttendanceStatusHadir
		a.IsLate = true
	} else {
		a.Status = models.AttendanceStatusHadir
		a.IsLate = false
	}

	a.Datetime = now
	a.IsLocationValid = true

	return s.absensi.Create(a)
}

func (s *KegiatanService) GetAttendanceBySerdik(serdikID uint64) ([]models.Absensi, *repository.AbsensiSummary, error) {
	list, err := s.absensi.GetBySerdikID(serdikID)
	if err != nil {
		return nil, nil, err
	}
	summary, err := s.absensi.GetSummaryBySerdik(serdikID)
	return list, summary, err
}

func (s *KegiatanService) GetAttendanceByKegiatan(kegiatanID uint64) ([]models.Absensi, error) {
	return s.absensi.GetByKegiatanID(kegiatanID)
}

// ---- Izin ----

func (s *KegiatanService) SubmitIzin(req *models.IzinRequest) error {
	return s.izin.Create(req)
}

func (s *KegiatanService) GetIzinBySerdik(serdikID int64) ([]models.IzinRequest, error) {
	return s.izin.GetBySerdikID(serdikID)
}

func (s *KegiatanService) GetPendingIzin() ([]models.IzinRequest, error) {
	return s.izin.GetPending()
}

func (s *KegiatanService) ReviewIzin(id int64, status string, reviewerID int64, reason *string) error {
	return s.izin.UpdateStatus(id, status, reviewerID, reason)
}
