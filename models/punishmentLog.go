package models

import (
	"time"
	"github.com/google/uuid"
)

// =========================
// PunishmentLog (model terpisah)
// =========================
type PunishmentLog struct {
    ID               int64      `gorm:"primaryKey"`
    UserID           int64      `gorm:"not null"`
    PunishmentItemID int64      `gorm:"not null"`
    FileID           *uuid.UUID `gorm:"type:uuid"`
    Qty              int        `gorm:"default:1"`
    Point            float64    `gorm:"type:numeric(8,2);not null"`
    ViolationDate    time.Time  `gorm:"not null"`
    Notes            *string
    CreatedBy        *int64
    CreatedAt        time.Time  `gorm:"autoCreateTime"`
    UpdatedAt        time.Time  `gorm:"autoUpdateTime"`

    // Relasi ke PunishmentItem
    PunishmentItem PunishmentItem `gorm:"foreignKey:PunishmentItemID"`
}