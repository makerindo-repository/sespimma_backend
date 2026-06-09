package dto

type PenilaianAkademikRequest struct {
	SerdikID            uint64  `json:"serdik_id" binding:"required"`
	AkademikComponentID uint64  `json:"akademik_component_id" binding:"required"`
	Nilai               float64 `json:"nilai" binding:"required"`
	Catatan             *string `json:"catatan"`
}