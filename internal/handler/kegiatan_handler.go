package handler

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"sespima_api/config"
	"sespima_api/internal/service"
	"sespima_api/models"
	"sespima_api/pkg/response"
)

type KegiatanHandler struct{ svc *service.KegiatanService }

func NewKegiatanHandler(svc *service.KegiatanService) *KegiatanHandler {
	return &KegiatanHandler{svc: svc}
}

func (h *KegiatanHandler) GetAll(c *gin.Context) {
	list, err := h.svc.GetAll()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *KegiatanHandler) GetRutin(c *gin.Context) {
	list, err := h.svc.GetRutin()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *KegiatanHandler) GetToday(c *gin.Context) {
	list, err := h.svc.GetToday()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *KegiatanHandler) Create(c *gin.Context) {
	var input models.Kegiatan
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	uid, exists := c.Get("user_id")
	if exists {
		input.CreatedID = int64(uid.(float64))
	}

	if err := h.svc.Create(&input); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, "kegiatan berhasil dibuat", input)
}

func (h *KegiatanHandler) Update(c *gin.Context) {
	var input models.Kegiatan
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	uid, exists := c.Get("user_id")
	if exists {
		input.CreatedID = int64(uid.(float64))
	}

	if err := h.svc.Update(c.Param("id"), &input); err != nil {
		response.NotFound(c, "kegiatan tidak ditemukan")
		return
	}
	response.OK(c, "kegiatan berhasil diperbarui", nil)
}

func (h *KegiatanHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "kegiatan berhasil dihapus", nil)
}

// ---- Absensi ----

func (h *KegiatanHandler) CheckIn(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists {
		response.Forbidden(c, "unauthorized")
		return
	}
	userID := uint64(uid.(float64))

	var serdik models.Serdik
	if err := config.DB.Where("user_id = ?", userID).First(&serdik).Error; err != nil {
		response.Forbidden(c, "hanya serdik yang dapat melakukan checkin")
		return
	}

	var input models.Absensi
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	input.UserID = &userID
	serdikID := uint64(serdik.ID)
	input.SerdikID = &serdikID

	if err := h.svc.CheckIn(&input); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, "absensi berhasil dicatat", input)
}

func (h *KegiatanHandler) GetAttendanceBySerdik(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("serdikId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid serdik_id")
		return
	}
	list, summary, err := h.svc.GetAttendanceRecap(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", gin.H{"records": list, "summary": summary})
}

func (h *KegiatanHandler) GetMyAttendance(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists {
		response.Forbidden(c, "unauthorized")
		return
	}
	userID := uint64(uid.(float64))

	var serdik models.Serdik
	if err := config.DB.Where("user_id = ?", userID).First(&serdik).Error; err != nil {
		response.Forbidden(c, "hanya serdik yang memiliki data absensi")
		return
	}

	list, summary, err := h.svc.GetAttendanceRecap(uint64(serdik.ID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", gin.H{"records": list, "summary": summary})
}

func (h *KegiatanHandler) GetAttendanceByKegiatan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("kegiatanId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid kegiatan_id")
		return
	}
	list, err := h.svc.GetAttendanceByKegiatan(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

// ---- Zones (mobile-compatible format) ----

// GetTodayZones returns today's kegiatan in the AttendanceZone shape the mobile consumes.
func (h *KegiatanHandler) GetTodayZones(c *gin.Context) {
	type ZoneResp struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		ActivityName string  `json:"activity_name"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		RadiusMeters float64 `json:"radius_meters"`
		StartTime    string  `json:"start_time"`
		EndTime      string  `json:"end_time"`
		Deadline      string           `json:"deadline"`
		CutoffTime    string           `json:"cutoff_time"`
		IsRoutine     bool             `json:"is_routine"`
		IsTraining    bool             `json:"is_training"`
		CreatedAt     string           `json:"created_at"`
		PolygonPoints *json.RawMessage `json:"polygon_points,omitempty"`
	}

	rows, err := config.DB.Raw(`
		SELECT
			id::text                                        AS id,
			nama,
			nama_lokasi,
			COALESCE(ST_Y(location::geometry), 0.0)        AS latitude,
			COALESCE(ST_X(location::geometry), 0.0)        AS longitude,
			COALESCE(radius, 100.0)                        AS radius,
			tanggal_mulai::text                            AS tanggal_mulai,
			waktu_mulai::text                              AS waktu_mulai,
			waktu_selesai::text                            AS waktu_selesai,
			batas_waktu_penugasan::text                    AS batas_waktu_penugasan,
			COALESCE(cutoff_time, batas_waktu_penugasan)::text AS cutoff_time,
			is_rutin,
			is_pelatihan,
			created_at::text                               AS created_at,
			polygon_points
		FROM kegiatan
		WHERE tanggal_mulai::date = CURRENT_DATE
		  AND location IS NOT NULL
	`).Rows()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	defer rows.Close()

	var zones []ZoneResp
	for rows.Next() {
		var z struct {
			ID, Nama, NamaLokasi, TanggalMulai, WaktuMulai, WaktuSelesai, BatasWaktu, Cutoff, CreatedAt string
			Lat, Lng, Radius                                                                            float64
			IsRutin, IsPelatihan                                                                        bool
			PolygonPoints                                                                               *json.RawMessage
		}
		if err := rows.Scan(&z.ID, &z.Nama, &z.NamaLokasi, &z.Lat, &z.Lng, &z.Radius,
			&z.TanggalMulai, &z.WaktuMulai, &z.WaktuSelesai, &z.BatasWaktu, &z.Cutoff,
			&z.IsRutin, &z.IsPelatihan, &z.CreatedAt, &z.PolygonPoints); err != nil {
			continue
		}
		base := z.TanggalMulai[:10] // "YYYY-MM-DD"
		zones = append(zones, ZoneResp{
			ID:            z.ID,
			Name:          z.NamaLokasi,
			ActivityName:  z.Nama,
			Latitude:      z.Lat,
			Longitude:     z.Lng,
			RadiusMeters:  z.Radius,
			StartTime:     fmt.Sprintf("%sT%s+07:00", base, z.WaktuMulai[:5]),
			EndTime:       fmt.Sprintf("%sT%s+07:00", base, z.WaktuSelesai[:5]),
			Deadline:      fmt.Sprintf("%sT%s+07:00", base, z.BatasWaktu[:5]),
			CutoffTime:    fmt.Sprintf("%sT%s+07:00", base, z.Cutoff[:5]),
			IsRoutine:     z.IsRutin,
			IsTraining:    z.IsPelatihan,
			CreatedAt:     z.CreatedAt,
			PolygonPoints: z.PolygonPoints,
		})
	}
	if zones == nil {
		zones = []ZoneResp{}
	}
	response.OK(c, "success", zones)
}

// ---- Izin ----

func (h *KegiatanHandler) SubmitIzin(c *gin.Context) {
	uid, _ := c.Get("user_id")
	userID := uint(uid.(float64))

	var serdik models.Serdik
	if err := config.DB.Where("user_id = ?", userID).First(&serdik).Error; err != nil {
		response.Forbidden(c, "data serdik tidak ditemukan untuk akun ini")
		return
	}

	var input models.IzinRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	input.SerdikID = int64(serdik.ID)

	if err := h.svc.SubmitIzin(&input); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, "permohonan izin berhasil diajukan", input)
}

func (h *KegiatanHandler) GetIzinBySerdik(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("serdikId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid serdik_id")
		return
	}
	list, err := h.svc.GetIzinBySerdik(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *KegiatanHandler) GetPendingIzin(c *gin.Context) {
	list, err := h.svc.GetPendingIzin()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

type reviewIzinReq struct {
	Status          string  `json:"status" binding:"required"`
	RejectionReason *string `json:"rejection_reason"`
}

func (h *KegiatanHandler) ReviewIzin(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	uid, _ := c.Get("user_id")
	reviewerID := int64(uid.(float64))

	var req reviewIzinReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.ReviewIzin(id, req.Status, reviewerID, req.RejectionReason); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "izin berhasil diproses", nil)
}
