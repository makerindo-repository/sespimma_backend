package models

import (
    "time"
)

type User struct {
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string  `gorm:"type:varchar(150);unique;not null" json:"email"`
	NRP          *string `gorm:"column:nrp;type:varchar(50);unique" json:"nrp"`
	NIP          *string `gorm:"column:nip;type:varchar(50);unique" json:"nip"`
	Password     string  `gorm:"type:varchar(255);not null" json:"-"`
	Role         string  `gorm:"type:user_role;default:students" json:"role"`
	IsActive     bool    `gorm:"default:true" json:"is_active"`
	IsFirstLogin bool    `gorm:"default:true" json:"is_first_login"`
	IsNakApproved bool   `gorm:"default:false" json:"is_nak_approved"`
	CurrentToken *string `gorm:"type:text" json:"-"`
	RefreshToken *string `gorm:"type:text" json:"-"`
	ResetToken   *string `gorm:"type:varchar(255)" json:"-"`
	FcmToken     *string `gorm:"type:varchar(500)" json:"-"`
	ProfilePhoto *string `gorm:"type:varchar(500)" json:"profile_photo,omitempty"`
	Location     *string `gorm:"type:geography(Geometry,4326)" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Pimpinan *Pimpinan `gorm:"foreignKey:UserID;references:ID" json:"pimpinan,omitempty"`
	Patun    *Patun    `gorm:"foreignKey:UserID;references:ID" json:"patun,omitempty"`
	Korsis   *Korsis   `gorm:"foreignKey:UserID;references:ID" json:"korsis,omitempty"`
	Serdik   *Serdik   `gorm:"foreignKey:UserID;references:ID" json:"serdik,omitempty"`
	Gadik    *Gadik    `gorm:"foreignKey:UserID;references:ID" json:"gadik,omitempty"`
}

func (User) TableName() string {
    return "users"
}

