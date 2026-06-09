package models

import (
    "time"
)

type User struct {
    ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Email        string    `gorm:"type:varchar(150);unique;not null" json:"email"`
    NRP          *string   `gorm:"type:varchar(50);unique" json:"nrp"`
    NIP          *string   `gorm:"type:varchar(50);unique" json:"nip"`
    Password     string    `gorm:"type:varchar(255);not null" json:"-"`
    Role         string    `gorm:"type:user_role;default:students" json:"role"`
    IsActive     bool      `gorm:"default:true" json:"is_active"`
    IsFirstLogin bool      `gorm:"default:true" json:"is_first_login"`
    CurrentToken *string   `gorm:"type:text" json:"current_token"`
    ResetToken   *string   `gorm:"type:varchar(255)" json:"reset_token,omitempty"`
    Location     *string   `gorm:"type:geography(Geometry,4326)" json:"location"`

    CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

    // One-to-One relation (pointer avoids recursive type cycle)
    Pimpinan     *Pimpinan `gorm:"foreignKey:UserID;references:ID" json:"pimpinan,omitempty"`
	Patun        *Patun    `gorm:"foreignKey:UserID;references:ID" json:"patun,omitempty"`
	Korsis       *Korsis   `gorm:"foreignKey:UserID;references:ID" json:"korsis,omitempty"`
    Serdik       *Serdik   `gorm:"foreignKey:UserID;references:ID" json:"serdik,omitempty"`
    Gadik         *Gadik     `gorm:"foreignKey:UserID;references:ID" json:"gadik,omitempty"`
}

func (User) TableName() string {
    return "users"
}

