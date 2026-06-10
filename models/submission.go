package models

import "time"

type Submission struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AssignmentID int64     `gorm:"not null;index" json:"assignment_id"`
	SerdikID     int64     `gorm:"not null;index" json:"serdik_id"`

	SubmittedAt *time.Time `json:"submitted_at"`
	FileURL     *string    `gorm:"type:varchar(1000)" json:"file_url"`
	FileName    *string    `gorm:"type:varchar(500)" json:"file_name"`

	IsGraded        bool     `gorm:"default:false" json:"is_graded"`
	IsRemedial      bool     `gorm:"default:false" json:"is_remedial"`
	Status          string   `gorm:"type:submission_status_enum;default:'pending'" json:"status"`
	CatatanPengajar *string  `gorm:"type:text" json:"catatan_pengajar"`
	NilaiAkhir      *float64 `json:"nilai_akhir"`

	ScoreMateri                 *float64 `json:"score_materi"`
	ScorePenulisan              *float64 `json:"score_penulisan"`
	ScorePaparan                *float64 `json:"score_paparan"`
	ScoreKeaktifan              *float64 `json:"score_keaktifan"`
	ScoreUjian                  *float64 `json:"score_ujian"`
	ScoreKeaktifanPerseorangan  *float64 `json:"score_keaktifan_perseorangan"`
	ScoreProdukPerseorangan     *float64 `json:"score_produk_perseorangan"`
	ScoreTataRuang              *float64 `json:"score_tata_ruang"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Assignment *Assignment `gorm:"foreignKey:AssignmentID;references:ID;constraint:OnDelete:CASCADE" json:"assignment,omitempty"`
	Serdik     *Serdik     `gorm:"foreignKey:SerdikID;references:ID;constraint:OnDelete:CASCADE" json:"serdik,omitempty"`
}

func (Submission) TableName() string { return "submissions" }
