package models

import "time"

// RewardPunishmentRule mirrors lib/core/constants/reward_punishment_data.dart.
// Aspect is one of the five Mental Kepribadian components.
type RewardPunishmentRule struct {
	Code        string    `gorm:"primaryKey;column:code" json:"code"`
	Type        string    `gorm:"column:type" json:"type"`     // REWARD | PUNISHMENT
	Aspect      string    `gorm:"column:aspect" json:"aspect"` // MORAL | DISIPLIN | KEPEMIMPINAN | PENGENDALIAN DIRI | PENAMPILAN
	Description string    `gorm:"column:description" json:"description"`
	Point       float64   `gorm:"column:point" json:"point"` // signed
	Note        *string   `gorm:"column:note" json:"note,omitempty"`
	IsActive    bool      `gorm:"column:is_active" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (RewardPunishmentRule) TableName() string { return "reward_punishment_rules" }

// RewardPunishmentRecord is a Maker-Checker reward/punishment assigned to a serdik.
type RewardPunishmentRecord struct {
	ID              int64      `gorm:"primaryKey" json:"id"`
	SerdikID        int64      `gorm:"column:serdik_id" json:"serdik_id"`
	RuleCode        string     `gorm:"column:rule_code" json:"rule_code"`
	Type            string     `gorm:"column:type" json:"type"`
	Aspect          string     `gorm:"column:aspect" json:"aspect"`
	Point           float64    `gorm:"column:point" json:"point"`
	Description     *string    `gorm:"column:description" json:"description"`
	Status          string     `gorm:"column:status;default:'pending'" json:"status"`
	CreatedBy       *int64     `gorm:"column:created_by" json:"created_by"`
	ApprovedBy      *int64     `gorm:"column:approved_by" json:"approved_by"`
	ReviewedAt      *time.Time `gorm:"column:reviewed_at" json:"reviewed_at"`
	RejectionReason *string    `gorm:"column:rejection_reason" json:"rejection_reason"`
	AttachmentPath  *string    `gorm:"column:attachment_path" json:"attachment_path"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	Serdik *Serdik               `gorm:"foreignKey:SerdikID;references:ID" json:"serdik,omitempty"`
	Rule   *RewardPunishmentRule `gorm:"foreignKey:RuleCode;references:Code" json:"rule,omitempty"`
}

func (RewardPunishmentRecord) TableName() string { return "reward_punishment_records" }
