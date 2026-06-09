package models

import (
    "time"
)

// =========================
// RewardCategory + RewardItem
// =========================
type RewardCategory struct {
    ID        int64          `gorm:"primaryKey"`
    Code      string         `gorm:"size:100;unique;not null"`
    Name      string         `gorm:"size:150;not null"`
    Bobot     float64        `gorm:"type:numeric(5,2);default:0"`
    SortOrder int            `gorm:"default:0"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`

    // Relasi ke RewardItem
    Items []RewardItem `gorm:"foreignKey:CategoryID"`
}

type RewardItem struct {
    ID         int64     `gorm:"primaryKey"`
    CategoryID int64     `gorm:"not null"`
    ItemNo     int       `gorm:"not null"`
    Kegiatan   string    `gorm:"type:text;not null"`
    Point      float64   `gorm:"type:numeric(6,2);default:0"`
    MaxCount   *int
    PeriodType string    `gorm:"type:reward_item_period_type;default:'unlimited'"`
    SpecialRule *string
    IsActive   bool      `gorm:"default:true"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime"`

    // Relasi balik ke RewardCategory
    Category RewardCategory `gorm:"foreignKey:CategoryID"`
}



