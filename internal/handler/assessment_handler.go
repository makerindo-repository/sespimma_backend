package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sespima_api/dto"
	"sespima_api/internal/service"
	"sespima_api/models"
	"sespima_api/pkg/response"
)

type AssessmentHandler struct{ svc *service.AssessmentService }

func NewAssessmentHandler(svc *service.AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{svc: svc}
}

// ---- Akademik ----

func (h *AssessmentHandler) GetAkademik(c *gin.Context) {
	tree, err := h.svc.GetAkademikTree()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", tree)
}

func (h *AssessmentHandler) GetAkademikBySerdik(c *gin.Context) {
	tree, err := h.svc.GetAkademikBySerdik(c.Param("serdikId"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", tree)
}

func (h *AssessmentHandler) CreateAkademik(c *gin.Context) {
	var req dto.PenilaianAkademikRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	p := &models.PenilaianAkademik{
		SerdikID:            req.SerdikID,
		AkademikComponentID: req.AkademikComponentID,
		Nilai:               req.Nilai,
		Catatan:             req.Catatan,
	}
	if err := h.svc.CreateAkademikPenilaian(p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "success", p)
}

func (h *AssessmentHandler) UpdateAkademik(c *gin.Context) {
	var req dto.PenilaianAkademikRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateAkademikPenilaian(c.Param("id"), req.Nilai, req.Catatan); err != nil {
		response.NotFound(c, "data tidak ditemukan")
		return
	}
	response.OK(c, "success", nil)
}

func (h *AssessmentHandler) DeleteAkademik(c *gin.Context) {
	if err := h.svc.DeleteAkademikPenilaian(c.Param("id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", nil)
}

// ---- Jasmani ----

func (h *AssessmentHandler) GetJasmani(c *gin.Context) {
	tree, err := h.svc.GetJasmaniTree()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", tree)
}

func (h *AssessmentHandler) GetJasmaniWithNilai(c *gin.Context) {
	comps, nilaiMap, err := h.svc.GetJasmaniWithNilai(c.Param("serdikId"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", gin.H{"components": comps, "nilai": nilaiMap})
}

func (h *AssessmentHandler) CreateJasmani(c *gin.Context) {
	var p models.PenilaianJasmani
	if err := c.ShouldBindJSON(&p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.CreateJasmaniPenilaian(&p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "success", p)
}

func (h *AssessmentHandler) UpdateJasmani(c *gin.Context) {
	var p models.PenilaianJasmani
	if err := c.ShouldBindJSON(&p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateJasmaniPenilaian(c.Param("id"), p.Nilai, p.Catatan); err != nil {
		response.NotFound(c, "data tidak ditemukan")
		return
	}
	response.OK(c, "success", nil)
}

func (h *AssessmentHandler) DeleteJasmani(c *gin.Context) {
	if err := h.svc.DeleteJasmaniPenilaian(c.Param("id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", nil)
}

// ---- Mental ----

func (h *AssessmentHandler) GetMental(c *gin.Context) {
	tree, err := h.svc.GetMentalTree()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", tree)
}

func (h *AssessmentHandler) GetMentalWithNilai(c *gin.Context) {
	comps, nilaiMap, err := h.svc.GetMentalWithNilai(c.Param("serdikId"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", gin.H{"components": comps, "nilai": nilaiMap})
}

func (h *AssessmentHandler) CreateMental(c *gin.Context) {
	var p models.PenilaianMental
	if err := c.ShouldBindJSON(&p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.CreateMentalPenilaian(&p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "success", p)
}

func (h *AssessmentHandler) UpdateMental(c *gin.Context) {
	var p models.PenilaianMental
	if err := c.ShouldBindJSON(&p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateMentalPenilaian(c.Param("id"), p.Nilai, p.Catatan); err != nil {
		response.NotFound(c, "data tidak ditemukan")
		return
	}
	response.OK(c, "success", nil)
}

func (h *AssessmentHandler) DeleteMental(c *gin.Context) {
	if err := h.svc.DeleteMentalPenilaian(c.Param("id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", nil)
}

// ---- Sociometry ----

func (h *AssessmentHandler) GetPeriods(c *gin.Context) {
	var pokjarID *int
	if v := c.Query("pokjar_id"); v != "" {
		id, _ := strconv.Atoi(v)
		pokjarID = &id
	}
	list, err := h.svc.GetPeriods(pokjarID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *AssessmentHandler) CreatePeriod(c *gin.Context) {
	var p models.SociometryPeriod
	if err := c.ShouldBindJSON(&p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	uid, _ := c.Get("user_id")
	createdBy := int64(uid.(float64))
	p.CreatedBy = &createdBy
	if err := h.svc.CreatePeriod(&p); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, "periode dibuat", p)
}

func (h *AssessmentHandler) ClosePeriod(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.ClosePeriod(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "periode ditutup", nil)
}

func (h *AssessmentHandler) SubmitEvaluation(c *gin.Context) {
	var e models.SociometryEvaluation
	if err := c.ShouldBindJSON(&e); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.SubmitEvaluation(&e); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "evaluasi berhasil disimpan", e)
}

// ---- Health ----

func (h *AssessmentHandler) GetHealthBySerdik(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("serdikId"), 10, 64)
	d, err := h.svc.GetHealthBySerdik(id)
	if err != nil {
		response.NotFound(c, "data kesehatan tidak ditemukan")
		return
	}
	response.OK(c, "success", d)
}

func (h *AssessmentHandler) UpsertHealthData(c *gin.Context) {
	var d models.SerdikHealthData
	if err := c.ShouldBindJSON(&d); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpsertHealthData(&d); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "data kesehatan berhasil disimpan", d)
}

func (h *AssessmentHandler) AddHealthRecord(c *gin.Context) {
	var r models.HealthRecord
	if err := c.ShouldBindJSON(&r); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	uid, _ := c.Get("user_id")
	r.MedisUserID = int64(uid.(float64))
	if err := h.svc.AddHealthRecord(&r); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, "catatan kesehatan berhasil ditambahkan", r)
}
