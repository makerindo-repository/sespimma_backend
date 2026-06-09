package models

import "time"

type AttendanceStatus string

const (
	AttendanceStatusHadir AttendanceStatus = "hadir"
	AttendanceStatusIzin  AttendanceStatus = "izin"
	AttendanceStatusSakit AttendanceStatus = "sakit"
	AttendanceStatusTK    AttendanceStatus = "tk"
)

type Absensi struct {
	ID         uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	SerdikID   uint64           `gorm:"not null;index" json:"serdik_id"`
	KegiatanID uint64           `gorm:"not null;index" json:"kegiatan_id"`

	Datetime time.Time        `gorm:"not null" json:"datetime"`
	Status   AttendanceStatus `gorm:"type:attendance_status_enum;not null" json:"status"`
	IsLate   bool             `gorm:"default:false" json:"is_late"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Serdik   Serdik   `gorm:"foreignKey:SerdikID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"serdik"`
	Kegiatan Kegiatan `gorm:"foreignKey:KegiatanID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"kegiatan"`
}

func (Absensi) TableName() string {
	return "absensi"
}