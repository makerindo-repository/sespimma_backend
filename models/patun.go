package models

import (
    "time"
)

type Patun struct {
    ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID            uint      `gorm:"unique;not null" json:"user_id"`

    Nama              string    `gorm:"type:varchar(150);not null" json:"nama"`
    Pangkat           *string   `gorm:"type:varchar(150)" json:"pangkat,omitempty"`
    NrpNip            string    `gorm:"type:varchar(50);unique;not null" json:"nrp_nip"`
    JabatanStruktural *string   `gorm:"type:varchar(150)" json:"jabatan_struktural,omitempty"`
    PeranPengasuhan   *string   `gorm:"type:varchar(150)" json:"peran_pengasuhan,omitempty"`
    Pokjar            *string   `gorm:"type:varchar(50)" json:"pokjar,omitempty"`

    CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

    // Relation to User
    User              *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	PokjarRel         *Pokjar   `gorm:"foreignKey:Pokjar;references:Name" json:"pokjar_rel,omitempty"`
}

// Optional: custom table name
func (Patun) TableName() string {
    return "patun"
}
