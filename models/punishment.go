package models

import (
    "time"

)

// =========================
// PunishmentCategory + PunishmentItem
// =========================
type PunishmentCategory struct {
    ID            int64     `gorm:"primaryKey"`
    Code          string    `gorm:"size:50;unique;not null"`
    Name          string    `gorm:"size:100;not null"`
    WeightPercent float64   `gorm:"type:numeric(5,2);default:0"`
    CreatedAt     time.Time `gorm:"autoCreateTime"`
    UpdatedAt     time.Time `gorm:"autoUpdateTime"`

    // Relasi ke PunishmentItem
    Items []PunishmentItem `gorm:"foreignKey:CategoryID"`
}

type PunishmentItem struct {
    ID         int64     `gorm:"primaryKey"`
    CategoryID int64     `gorm:"not null"`
    ItemNo     int       `gorm:"not null"`
    Activity   string    `gorm:"type:text;not null"`
    Point      float64   `gorm:"type:numeric(5,2);not null"`
    IsActive   bool      `gorm:"default:true"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime"`

    // Relasi balik ke PunishmentCategory
    Category PunishmentCategory `gorm:"foreignKey:CategoryID"`
}

