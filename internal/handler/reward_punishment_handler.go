package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"sespima_api/config"
	"sespima_api/models"
	"sespima_api/pkg/response"
)

// RewardPunishmentHandler drives the aspect-based (Mental Kepribadian)
// reward/punishment catalog, live serdik search, and Maker-Checker records.
type RewardPunishmentHandler struct{}

func NewRewardPunishmentHandler() *RewardPunishmentHandler { return &RewardPunishmentHandler{} }

// GetRules returns the catalog, optionally filtered by ?type= and ?aspect=.
func (h *RewardPunishmentHandler) GetRules(c *gin.Context) {
	q := config.DB.Model(&models.RewardPunishmentRule{}).Where("is_active = TRUE")
	if t := c.Query("type"); t != "" {
		q = q.Where("type = ?", strings.ToUpper(t))
	}
	if a := c.Query("aspect"); a != "" {
		q = q.Where("aspect = ?", strings.ToUpper(a))
	}
	var rules []models.RewardPunishmentRule
	if err := q.Order("type, aspect, code").Find(&rules).Error; err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", rules)
}

// SearchSerdik is the live DB lookup makers use when assigning a reward/penalty.
// Matches no_serdik / nrp / nama_lengkap (case-insensitive).
func (h *RewardPunishmentHandler) SearchSerdik(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if len(q) < 2 {
		response.BadRequest(c, "query minimal 2 karakter")
		return
	}

	type serdikHit struct {
		ID            int64   `json:"id"`
		NoSerdik      *string `json:"no_serdik"`
		NRP           *string `json:"nrp"`
		NamaLengkap   *string `json:"nama_lengkap"`
		Pangkat       *string `json:"pangkat"`
		KelompokKelas *string `json:"kelompok_kelas"`
	}
	like := "%" + q + "%"
	var hits []serdikHit
	if err := config.DB.Raw(`
		SELECT id, no_serdik, nrp, nama_lengkap, pangkat, kelompok_kelas
		FROM serdik
		WHERE no_serdik ILIKE ? OR nrp ILIKE ? OR nama_lengkap ILIKE ?
		ORDER BY nama_lengkap
		LIMIT 30`, like, like, like).Scan(&hits).Error; err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", hits)
}

type createRPReq struct {
	SerdikID       int64   `json:"serdik_id" binding:"required"`
	RuleCode       string  `json:"rule_code" binding:"required"`
	Description    *string `json:"description"`
	AttachmentPath *string `json:"attachment_path"`
}

// CreateRecord is the MAKER step: patun/gadik assign a reward/penalty.
// Always persisted as 'pending'; the point/type/aspect are snapshotted from
// the catalog rule so later edits to the rule never mutate history.
func (h *RewardPunishmentHandler) CreateRecord(c *gin.Context) {
	makerID, ok := userIDFromCtx(c)
	if !ok {
		response.Unauthorized(c, "sesi tidak valid")
		return
	}

	var req createRPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var rule models.RewardPunishmentRule
	if err := config.DB.Where("code = ? AND is_active = TRUE", req.RuleCode).First(&rule).Error; err != nil {
		response.BadRequest(c, "kode reward/punishment tidak ditemukan")
		return
	}

	// Verify target serdik exists.
	var count int64
	config.DB.Model(&models.Serdik{}).Where("id = ?", req.SerdikID).Count(&count)
	if count == 0 {
		response.NotFound(c, "serdik tidak ditemukan")
		return
	}

	role, _ := c.Get("role")
	status := "pending"
	msg := "pengajuan berhasil dikirim, menunggu persetujuan Korsis"
	
	if role == "korsis" || role == "operator" {
		status = "approved"
		msg = "reward/punishment berhasil ditambahkan"
	}

	mid := int64(makerID)
	rec := models.RewardPunishmentRecord{
		SerdikID:       req.SerdikID,
		RuleCode:       rule.Code,
		Type:           rule.Type,
		Aspect:         rule.Aspect,
		Point:          rule.Point,
		Description:    req.Description,
		AttachmentPath: req.AttachmentPath,
		Status:         status,
		CreatedBy:      &mid,
	}
	if status == "approved" {
		rec.ApprovedBy = &mid
		now := time.Now()
		rec.ReviewedAt = &now
	}

	if err := config.DB.Create(&rec).Error; err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, msg, rec)
}

// GetPending lists all pending records for the Korsis approval inbox.
func (h *RewardPunishmentHandler) GetPending(c *gin.Context) {
	var recs []models.RewardPunishmentRecord
	if err := config.DB.Where("status = 'pending'").
		Preload("Serdik").Preload("Rule").
		Order("created_at DESC").Find(&recs).Error; err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", recs)
}

// GetBySerdik lists records for a given serdik (any role with the route guard).
func (h *RewardPunishmentHandler) GetBySerdik(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("serdikId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid serdik_id")
		return
	}
	var recs []models.RewardPunishmentRecord
	if err := config.DB.Where("serdik_id = ?", id).
		Preload("Rule").Order("created_at DESC").Find(&recs).Error; err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", recs)
}

// GetMyRecords lists the authenticated serdik's own approved/pending records.
func (h *RewardPunishmentHandler) GetMyRecords(c *gin.Context) {
	userID, ok := userIDFromCtx(c)
	if !ok {
		response.Unauthorized(c, "sesi tidak valid")
		return
	}
	var serdik models.Serdik
	if err := config.DB.Where("user_id = ?", userID).First(&serdik).Error; err != nil {
		response.NotFound(c, "data serdik tidak ditemukan")
		return
	}
	var recs []models.RewardPunishmentRecord
	if err := config.DB.Where("serdik_id = ?", serdik.ID).
		Preload("Rule").Order("created_at DESC").Find(&recs).Error; err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", recs)
}

// ApproveRecord approves a pending RP record
func (h *RewardPunishmentHandler) ApproveRecord(c *gin.Context) {
	checkerID, ok := userIDFromCtx(c)
	if !ok {
		response.Unauthorized(c, "sesi tidak valid")
		return
	}

	id := c.Param("id")
	var rec models.RewardPunishmentRecord
	if err := config.DB.Where("id = ? AND status = 'pending'", id).First(&rec).Error; err != nil {
		response.NotFound(c, "data pending tidak ditemukan")
		return
	}

	cid := int64(checkerID)
	now := time.Now()
	rec.Status = "approved"
	rec.ApprovedBy = &cid
	rec.ReviewedAt = &now

	if err := config.DB.Save(&rec).Error; err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "pengajuan berhasil disetujui", rec)
}

// RejectRecord rejects a pending RP record
func (h *RewardPunishmentHandler) RejectRecord(c *gin.Context) {
	checkerID, ok := userIDFromCtx(c)
	if !ok {
		response.Unauthorized(c, "sesi tidak valid")
		return
	}

	id := c.Param("id")
	var rec models.RewardPunishmentRecord
	if err := config.DB.Where("id = ? AND status = 'pending'", id).First(&rec).Error; err != nil {
		response.NotFound(c, "data pending tidak ditemukan")
		return
	}

	var req struct {
		Reason string `json:"rejection_reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cid := int64(checkerID)
	now := time.Now()
	rec.Status = "rejected"
	rec.ApprovedBy = &cid
	rec.ReviewedAt = &now
	rec.RejectionReason = &req.Reason

	if err := config.DB.Save(&rec).Error; err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "pengajuan ditolak", rec)
}
