package models

import "time"

type SociometryPeriod struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PokjarID   int       `gorm:"not null;index" json:"pokjar_id"`
	PeriodType string    `gorm:"type:sociometry_period_type;not null" json:"period_type"`
	StartDate  time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate    time.Time `gorm:"type:date;not null" json:"end_date"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedBy  *int64    `json:"created_by"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Pokjar    *Pokjar `gorm:"foreignKey:PokjarID;references:ID;constraint:OnDelete:RESTRICT" json:"pokjar,omitempty"`
	CreatedByUser *User `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnDelete:SET NULL" json:"created_by_user,omitempty"`
}

func (SociometryPeriod) TableName() string { return "sociometry_periods" }

type SociometryEvaluation struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PeriodID           int64     `gorm:"not null;index" json:"period_id"`
	EvaluatorSerdikID  int64     `gorm:"not null;index" json:"evaluator_serdik_id"`
	EvaluatedSerdikID  int64     `gorm:"not null;index" json:"evaluated_serdik_id"`
	Score              float64   `gorm:"not null;default:0" json:"score"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Period    *SociometryPeriod `gorm:"foreignKey:PeriodID;references:ID;constraint:OnDelete:CASCADE" json:"period,omitempty"`
	Evaluator *Serdik           `gorm:"foreignKey:EvaluatorSerdikID;references:ID;constraint:OnDelete:CASCADE" json:"evaluator,omitempty"`
	Evaluated *Serdik           `gorm:"foreignKey:EvaluatedSerdikID;references:ID;constraint:OnDelete:CASCADE" json:"evaluated,omitempty"`
}

func (SociometryEvaluation) TableName() string { return "sociometry_evaluations" }
