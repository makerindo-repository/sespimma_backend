package models

import (
	"time"
)

type Pimpinan struct {
    ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID            uint      `gorm:"unique;not null" json:"user_id"`
    Nama              string    `gorm:"type:varchar(150);not null" json:"nama"`
    Pangkat           string    `gorm:"type:varchar(100);not null" json:"pangkat"`
    NrpNip            string    `gorm:"type:varchar(50);unique;not null" json:"nrp_nip"`
    JabatanStruktural string    `gorm:"type:varchar(255);not null" json:"jabatan_struktural"`
    PeranPengasuhan   string    `gorm:"type:varchar(100);not null" json:"peran_pengasuhan"`

    CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

    // Back-reference (pointer avoids cycle)
    User              *User     `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (Pimpinan) TableName() string {
    return "pimpinan"
}
