package models

import (
    "time"
)

type Pokjar struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `gorm:"type:varchar(255);unique;not null" json:"name"`
    Grade     string    `gorm:"type:varchar(50);not null" json:"grade"`

    PatunID   *uint     `gorm:"index" json:"patun_id,omitempty"`
    Patun     *Patun    `gorm:"foreignKey:PatunID;references:ID;constraint:OnDelete:SET NULL" json:"patun,omitempty"`

    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Optional: custom table name
func (Pokjar) TableName() string {
    return "pokjar"
}
