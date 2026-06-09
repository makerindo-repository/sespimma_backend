package models

import (
    "time"
)

type Jasmani struct {
    ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
    ParentID    *uint      `gorm:"type:bigint" json:"parent_id"`
	
    Code        string     `gorm:"type:varchar(50);unique" json:"code"`
    Name        string     `gorm:"type:varchar(255);not null" json:"name"`
    Weight      float64    `gorm:"type:numeric(5,2);not null" json:"weight"`
    AgeGroup    *string    `gorm:"type:jasmani_age_group" json:"age_group"`
    Description *string    `gorm:"type:text" json:"description"`
	
    CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
    Parent      *Jasmani   `gorm:"foreignKey:ParentID;references:ID" json:"parent,omitempty"`
	Children    []Jasmani  `gorm:"foreignKey:ParentID;references:ID" json:"children"`
}

// Optional: custom table name
func (Jasmani) TableName() string {
    return "jasmani_components"
}
