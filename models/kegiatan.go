package models

import (
	"encoding/json"
	"time"
)

type Kegiatan struct {
	ID                  int64     `gorm:"primaryKey"`
	CreatedID           int64     `gorm:"not null"`
	Nama                string    `gorm:"size:255;not null"`
	NamaLokasi          string    `gorm:"size:255;not null"`
	TanggalMulai        time.Time `gorm:"not null"`
	WaktuMulai          string    `gorm:"type:time;not null"`
	WaktuSelesai        string    `gorm:"type:time;not null"`
	BatasWaktuPenugasan string    `gorm:"type:time;not null"`
	CutoffTime          *string   `gorm:"type:time" json:"cutoff_time,omitempty"`
	IsRutin             bool      `gorm:"default:false"`
	IsPelatihan         bool      `gorm:"default:false"`
	QRCode              bool      `gorm:"default:false"`
	Radius              *float64

	Location     string           `gorm:"type:geography(Geometry,4326)" json:"location"`
	PolygonPoints *json.RawMessage `gorm:"type:jsonb" json:"polygon_points,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (k *Kegiatan) TableName() string {
    return "kegiatan"
}
