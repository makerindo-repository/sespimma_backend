package models

import "time"

type AttendanceStatus string

const (
	AttendanceStatusHadir AttendanceStatus = "hadir"
	AttendanceStatusIzin  AttendanceStatus = "izin"
	AttendanceStatusSakit AttendanceStatus = "sakit"
	AttendanceStatusTK    AttendanceStatus = "tk"
)

type AttendanceType string

const (
	AttendanceTypeCheckin    AttendanceType = "checkin"
	AttendanceTypeCheckout   AttendanceType = "checkout"
	AttendanceTypeAbsent     AttendanceType = "absent"
	AttendanceTypePermission AttendanceType = "permission"
)

type Absensi struct {
	ID         uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	SerdikID   *uint64          `gorm:"index" json:"serdik_id"`
	KegiatanID *uint64          `gorm:"index" json:"kegiatan_id"`
	UserID     *uint64          `gorm:"index" json:"user_id"`

	Datetime        time.Time        `gorm:"not null" json:"datetime"`
	Status          AttendanceStatus `gorm:"type:attendance_status_enum;not null" json:"status"`
	Type            AttendanceType   `gorm:"type:varchar(20);default:'checkin'" json:"type"`
	Method          string           `gorm:"type:varchar(20);default:'gps'" json:"method"`
	IsLate          bool             `gorm:"default:false" json:"is_late"`
	IsLocationValid bool             `gorm:"default:false" json:"is_location_valid"`
	Latitude        *float64         `json:"latitude"`
	Longitude       *float64         `json:"longitude"`
	Note            *string          `gorm:"type:text" json:"note"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Serdik   *Serdik   `gorm:"foreignKey:SerdikID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"serdik,omitempty"`
	Kegiatan *Kegiatan `gorm:"foreignKey:KegiatanID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"kegiatan,omitempty"`
	User     *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"user,omitempty"`
}

func (Absensi) TableName() string {
	return "absensi"
}