package models

import "time"

type PenilaianMental struct {
	ID                uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	SerdikID          uint64  `gorm:"not null;uniqueIndex:uq_serdik_component" json:"serdik_id"`
	MentalComponentID uint64  `gorm:"not null;uniqueIndex:uq_serdik_component" json:"mental_component_id"`
	Nilai             float64 `gorm:"type:double precision;not null" json:"nilai"`
	Catatan           *string `gorm:"type:text" json:"catatan"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Serdik          Serdik `gorm:"foreignKey:SerdikID;references:ID;constraint:OnDelete:CASCADE" json:"serdik"`
	MentalComponent Mental `gorm:"foreignKey:MentalComponentID;references:ID;constraint:OnDelete:CASCADE" json:"mental_component"`
}

func (PenilaianMental) TableName() string {
	return "penilaian_mental"
}
