package models

import (
    "time"
)

type Korsis struct {
    ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID            uint      `gorm:"unique;not null" json:"user_id"`

    Nama              string    `gorm:"type:varchar(150);not null" json:"nama"`
    Pangkat           string    `gorm:"type:varchar(100);not null" json:"pangkat"`
    NrpNip            *string   `gorm:"type:varchar(50);unique" json:"nrp_nip"`
    JabatanStruktural *string   `gorm:"type:varchar(150)" json:"jabatan_struktural,omitempty"`
    PeranPengasuhan   *string   `gorm:"type:varchar(150)" json:"peran_pengasuhan,omitempty"`

    CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

    // Relation back to User
    User              User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
}

func (Korsis) TableName() string {
    return "korsis"
}
