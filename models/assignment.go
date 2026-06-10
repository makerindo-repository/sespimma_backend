package models

import "time"

type Assignment struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedBy      int64     `gorm:"not null" json:"created_by"`
	Judul          string    `gorm:"type:varchar(500);not null" json:"judul"`
	Deskripsi      *string   `gorm:"type:text" json:"deskripsi"`
	JenisTugas     *string   `gorm:"type:varchar(100)" json:"jenis_tugas"`
	TurunanTugas   *string   `gorm:"type:varchar(100)" json:"turunan_tugas"`
	Mapel          *string   `gorm:"type:varchar(255)" json:"mapel"`
	Deadline       time.Time `gorm:"not null" json:"deadline"`
	TargetPokjarID *int      `gorm:"index" json:"target_pokjar_id"`
	Instruksi      *string   `gorm:"type:text" json:"instruksi"`
	Status         string    `gorm:"type:assignment_status_enum;default:'active'" json:"status"`
	FileName       *string   `gorm:"type:varchar(500)" json:"file_name"`
	FileURL        *string   `gorm:"type:varchar(1000)" json:"file_url"`
	IsRemedial     bool      `gorm:"default:false" json:"is_remedial"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Creator *User   `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnDelete:RESTRICT" json:"creator,omitempty"`
	Pokjar  *Pokjar `gorm:"foreignKey:TargetPokjarID;references:ID;constraint:OnDelete:SET NULL" json:"pokjar,omitempty"`
}

func (Assignment) TableName() string { return "assignments" }
