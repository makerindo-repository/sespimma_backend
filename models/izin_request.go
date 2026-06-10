package models

import "time"

type IzinRequest struct {
	ID                 int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	SerdikID           int64      `gorm:"not null;index" json:"serdik_id"`
	KegiatanID         *int64     `gorm:"index" json:"kegiatan_id"`
	StartTime          time.Time  `gorm:"not null" json:"start_time"`
	EndTime            time.Time  `gorm:"not null" json:"end_time"`
	Description        string     `gorm:"type:text;not null" json:"description"`
	AttachmentPath     *string    `gorm:"type:varchar(1000)" json:"attachment_path"`
	AttachmentFileName *string    `gorm:"type:varchar(500)" json:"attachment_file_name"`
	ReviewedBy         *int64     `json:"reviewed_by"`
	ReviewedAt         *time.Time `json:"reviewed_at"`
	RejectionReason    *string    `gorm:"type:text" json:"rejection_reason"`
	Status             string     `gorm:"type:izin_status_enum;default:'pending'" json:"status"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Serdik   *Serdik   `gorm:"foreignKey:SerdikID;references:ID;constraint:OnDelete:CASCADE" json:"serdik,omitempty"`
	Kegiatan *Kegiatan `gorm:"foreignKey:KegiatanID;references:ID;constraint:OnDelete:SET NULL" json:"kegiatan,omitempty"`
	Reviewer *User     `gorm:"foreignKey:ReviewedBy;references:ID;constraint:OnDelete:SET NULL" json:"reviewer,omitempty"`
}

func (IzinRequest) TableName() string { return "izin_requests" }
