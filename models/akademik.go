package models

import (
    "time"
)

type Akademik struct {
    ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
    ParentID    *uint      `gorm:"type:bigint" json:"parent_id"`
    Code        string     `gorm:"type:varchar(50);unique" json:"code"`
    Name        string     `gorm:"type:varchar(255);not null" json:"name"`
    Weight      float64    `gorm:"type:numeric(5,2);not null" json:"weight"`
    Level       int16      `gorm:"type:smallint;default:1" json:"level"`
    Description *string    `gorm:"type:text" json:"description"`

    CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

    // Self-referencing relationship
    Parent      *Akademik  `gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:SET NULL" json:"parent"`
    Children    []Akademik `gorm:"foreignKey:ParentID;references:ID" json:"children"`
}

// Optional: custom table name
func (Akademik) TableName() string {
    return "akademik_component"
}
