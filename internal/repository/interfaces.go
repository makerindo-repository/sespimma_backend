package repository

import (
	"time"

	"sespima_api/models"
)

type UserRepository interface {
	FindByNrpOrNip(nrpNip string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	Save(user *models.User) error
	UpdateLocation(userID interface{}, lat, lng float64) error
	GetAllWithProfiles() ([]models.User, error)
	GetUserLocation(id string) (*UserLocationResult, error)
	GetAllUserLocations() ([]UserLocationResult, error)
}

type UserLocationResult struct {
	ID        uint     `json:"id"`
	Email     string   `json:"email"`
	Role      string   `json:"role,omitempty"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type LocationLogRepository interface {
	Create(log *models.UserLocationLog) error
	GetLatestByUser(userID int64) (*models.UserLocationLog, error)
}

type KegiatanRepository interface {
	GetAll() ([]models.Kegiatan, error)
	GetRutin() ([]models.Kegiatan, error)
	GetToday() ([]models.Kegiatan, error)
	GetByID(id string) (*models.Kegiatan, error)
	GetDistance(kegiatanID uint64, lat, lng float64) (float64, error)
	Create(k *models.Kegiatan) error
	Update(k *models.Kegiatan) error
	Delete(id string) error
}

type AbsensiRepository interface {
	Create(a *models.Absensi) error
	GetBySerdikID(serdikID uint64) ([]models.Absensi, error)
	GetByUserID(userID uint64) ([]models.Absensi, error)
	GetByKegiatanID(kegiatanID uint64) ([]models.Absensi, error)
	GetSummaryBySerdik(serdikID uint64) (*AbsensiSummary, error)
	// GetRecapBySerdik derives the strict per-activity attendance recap
	// (hadir/telat/izin/tk) by cross-referencing check-ins, approved izin,
	// and missed zones. Returns derived records plus an accurate summary.
	GetRecapBySerdik(serdikID uint64) ([]RecapRecord, *AbsensiSummary, error)
}

type AbsensiSummary struct {
	PresentCount    int `json:"present_count"`
	LateCount       int `json:"late_count"`
	PermissionCount int `json:"permission_count"`
	AbsentCount     int `json:"absent_count"`
}

// RecapKegiatan is the zone/activity context embedded in a recap record.
type RecapKegiatan struct {
	Nama       string `json:"nama"`
	NamaLokasi string `json:"nama_lokasi"`
}

// RecapRecord is one derived attendance row for a single zone/activity.
type RecapRecord struct {
	ID         *uint64       `json:"id"`          // absensi id, null if no check-in
	KegiatanID uint64        `json:"kegiatan_id"` // the zone id
	Datetime   time.Time     `json:"datetime"`
	Status     string        `json:"status"` // hadir | telat | izin | tk
	Method     string        `json:"method"` // gps | qr_code | -
	IsLate     bool          `json:"is_late"`
	Kegiatan   RecapKegiatan `json:"kegiatan"`
}

type AkademikRepository interface {
	GetAllComponents() ([]models.Akademik, error)
	GetPenilaianBySerdik(serdikID string) ([]models.PenilaianAkademik, error)
	CreatePenilaian(p *models.PenilaianAkademik) error
	UpdatePenilaian(id string, nilai float64, catatan *string) error
	DeletePenilaian(id string) error
}

type JasmaniRepository interface {
	GetAllComponents() ([]models.Jasmani, error)
	GetPenilaianBySerdik(serdikID string) ([]models.PenilaianJasmani, error)
	CreatePenilaian(p *models.PenilaianJasmani) error
	UpdatePenilaian(id string, nilai float64, catatan *string) error
	DeletePenilaian(id string) error
}

type MentalRepository interface {
	GetAllComponents() ([]models.Mental, error)
	GetPenilaianBySerdik(serdikID string) ([]models.PenilaianMental, error)
	CreatePenilaian(p *models.PenilaianMental) error
	UpdatePenilaian(id string, nilai float64, catatan *string) error
	DeletePenilaian(id string) error
}

type RewardRepository interface {
	GetAll() ([]models.UserReward, error)
	GetPending() ([]models.UserReward, error)
	Create(r *models.UserReward) error
	Update(id string, input *models.UserReward) error
	Delete(id string) error
	Approve(id string, approverID int64) error
	Reject(id string, approverID int64, reason string) error
}

type PunishmentRepository interface {
	GetAll() ([]models.PunishmentLog, error)
	GetByUserID(userID string) ([]models.PunishmentLog, error)
	GetPending() ([]models.PunishmentLog, error)
	Create(p *models.PunishmentLog) error
	Update(id string, input *models.PunishmentLog) error
	Delete(id string) error
	Approve(id string, approverID int64) error
	Reject(id string, approverID int64, reason string) error
}

type IzinRepository interface {
	Create(r *models.IzinRequest) error
	GetBySerdikID(serdikID int64) ([]models.IzinRequest, error)
	GetPending() ([]models.IzinRequest, error)
	UpdateStatus(id int64, status string, reviewerID int64, reason *string) error
}

type AssignmentRepository interface {
	GetAll(pokjarID *int) ([]models.Assignment, error)
	GetByID(id int64) (*models.Assignment, error)
	Create(a *models.Assignment) error
	Update(a *models.Assignment) error
	Delete(id int64) error
	GetSubmissions(assignmentID int64) ([]models.Submission, error)
	CreateSubmission(s *models.Submission) error
	UpdateSubmission(s *models.Submission) error
}

type SociometryRepository interface {
	GetPeriods(pokjarID *int) ([]models.SociometryPeriod, error)
	CreatePeriod(p *models.SociometryPeriod) error
	ClosePeriod(id int64) error
	CreateEvaluation(e *models.SociometryEvaluation) error
	GetEvaluationsByPeriod(periodID int64) ([]models.SociometryEvaluation, error)
	GetAverageScore(serdikID int64, periodType string) (float64, error)
}

type HealthRepository interface {
	GetBySerdikID(serdikID int64) (*models.SerdikHealthData, error)
	UpsertHealthData(d *models.SerdikHealthData) error
	CreateRecord(r *models.HealthRecord) error
	GetRecords(healthDataID int64) ([]models.HealthRecord, error)
}
