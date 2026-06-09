package models

import "time"

type PenilaianJasmani struct {
	ID                 uint64   `gorm:"primaryKey;autoIncrement" json:"id"`
	SerdikID           uint64   `gorm:"not null;index" json:"serdik_id"`
	JasmaniComponentID uint64   `gorm:"not null;index" json:"jasmani_component_id"`
	Nilai              float64  `gorm:"type:double precision;not null" json:"nilai"`
	Catatan            *string  `gorm:"type:text" json:"catatan"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Serdik          Serdik           `gorm:"foreignKey:SerdikID;references:ID;constraint:OnDelete:CASCADE" json:"serdik"`
	JasmaniComponent Jasmani `gorm:"foreignKey:JasmaniComponentID;references:ID;constraint:OnDelete:CASCADE" json:"jasmani_component"`
}

func (PenilaianJasmani) TableName() string {
	return "penilaian_jasmani"
}