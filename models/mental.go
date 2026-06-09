package models

import (
    "time"
)

// MentalIndicatorType corresponds to the ENUM mental_indicator_type in PostgreSQL
type MentalIndicatorType string

const (
    IndicatorComponent  MentalIndicatorType = "component"
    IndicatorReward     MentalIndicatorType = "reward"
    IndicatorPunishment MentalIndicatorType = "punishment"
)

// Mental represents the mental_components table
type Mental struct {
    ID           int64              `gorm:"primaryKey;autoIncrement"`
    ParentID     *int64             `gorm:"index"`
    Code         string             `gorm:"size:50;unique;not null"`
    Name         string             `gorm:"size:255;not null"`
    Weight       float64            `gorm:"type:numeric(5,2);not null"`
    IndicatorType MentalIndicatorType `gorm:"type:mental_indicator_type;default:'component'"`
    Point        *float64           `gorm:"type:numeric(6,2)"`
    Description  string             `gorm:"type:text"`
    CreatedAt    time.Time          `gorm:"default:current_timestamp"`
    UpdatedAt    time.Time          `gorm:"default:current_timestamp"`

    // Self-referencing relation
    Parent   *Mental   `gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:SET NULL"`
    Children []Mental  `gorm:"foreignKey:ParentID"`
}

func (Mental) TableName() string {
	return "mental_components"
}