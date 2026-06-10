package models

import "time"

type UserLocationLog struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"not null;index" json:"user_id"`
	Latitude   float64   `gorm:"not null" json:"latitude"`
	Longitude  float64   `gorm:"not null" json:"longitude"`
	Accuracy   *float64  `json:"accuracy"`
	RecordedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"recorded_at"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (UserLocationLog) TableName() string { return "user_location_logs" }
