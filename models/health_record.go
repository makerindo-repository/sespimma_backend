package models

import "time"

type SerdikHealthData struct {
	ID             int64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SerdikID       int64    `gorm:"not null;uniqueIndex" json:"serdik_id"`
	NilaiA         *float64 `json:"nilai_a"`
	CatatanDokterA *string  `gorm:"type:text" json:"catatan_dokter_a"`
	NilaiB         *float64 `json:"nilai_b"`
	CatatanDokterB *string  `gorm:"type:text" json:"catatan_dokter_b"`
	BaseNilaiC     int      `gorm:"default:80" json:"base_nilai_c"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Serdik  *Serdik       `gorm:"foreignKey:SerdikID;references:ID;constraint:OnDelete:CASCADE" json:"serdik,omitempty"`
	Records []HealthRecord `gorm:"foreignKey:SerdikHealthDataID;references:ID" json:"records,omitempty"`
}

func (SerdikHealthData) TableName() string { return "serdik_health_data" }

type HealthRecord struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SerdikHealthDataID  int64     `gorm:"not null;index" json:"serdik_health_data_id"`
	MedisUserID         int64     `gorm:"not null;index" json:"medis_user_id"`
	Type                string    `gorm:"type:varchar(100);not null" json:"type"`
	Description         string    `gorm:"type:text;not null" json:"description"`
	PhotoPath           *string   `gorm:"type:varchar(1000)" json:"photo_path"`
	MinusPoints         int       `gorm:"default:0" json:"minus_points"`
	RecordedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"recorded_at"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`

	HealthData *SerdikHealthData `gorm:"foreignKey:SerdikHealthDataID;references:ID;constraint:OnDelete:CASCADE" json:"health_data,omitempty"`
	MedisUser  *User             `gorm:"foreignKey:MedisUserID;references:ID;constraint:OnDelete:RESTRICT" json:"medis_user,omitempty"`
}

func (HealthRecord) TableName() string { return "health_records" }
