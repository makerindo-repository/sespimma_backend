package models

import (
	"github.com/google/uuid"
	"time"
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

	// Maker-Checker workflow (added in migration 033)
	Status          string     `gorm:"type:punishment_log_status;default:'pending'" json:"status"`
	ApprovedBy      *int64     `json:"approved_by"`
	ReviewedAt      *time.Time `json:"reviewed_at"`
	RejectionReason *string    `gorm:"type:text" json:"rejection_reason"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relasi ke PunishmentItem
	PunishmentItem PunishmentItem `gorm:"foreignKey:PunishmentItemID"`
}
