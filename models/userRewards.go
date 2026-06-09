package models

import (
	"time"
	"github.com/google/uuid"
)

// =========================
// UserReward (model terpisah)
// =========================
type UserReward struct {
    ID           int64      `gorm:"primaryKey"`
    UserID       int64      `gorm:"not null"`
    RewardItemID int64      `gorm:"not null"`
    FileID       *uuid.UUID `gorm:"type:uuid"`
    Qty          int        `gorm:"default:1"`
    Point        float64    `gorm:"type:numeric(8,2);not null"`
    RewardDate   time.Time  `gorm:"not null"`
    ApprovedBy   *int64
    CreatedBy    *int64
    Notes        *string
    Status       string     `gorm:"type:user_reward_status;default:'approved'"`
    CreatedAt    time.Time  `gorm:"autoCreateTime"`
    UpdatedAt    time.Time  `gorm:"autoUpdateTime"`

    // Relasi ke RewardItem
    RewardItem RewardItem `gorm:"foreignKey:RewardItemID"`
}