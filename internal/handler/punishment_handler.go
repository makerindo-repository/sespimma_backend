package handler

import (
	"sespima_api/internal/service"
	"sespima_api/models"
	"sespima_api/pkg/response"

	"github.com/gin-gonic/gin"
)

type PunishmentHandler struct{ svc *service.PunishmentService }

func NewPunishmentHandler(svc *service.PunishmentService) *PunishmentHandler {
	return &PunishmentHandler{svc: svc}
}

func (h *PunishmentHandler) GetAll(c *gin.Context) {
	list, err := h.svc.GetAll()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *PunishmentHandler) GetByUserID(c *gin.Context) {
	list, err := h.svc.GetByUserID(c.Param("user_id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *PunishmentHandler) GetPending(c *gin.Context) {
	list, err := h.svc.GetPending()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *PunishmentHandler) Approve(c *gin.Context) {
	uid, _ := c.Get("user_id")
	approverID := int64(uid.(float64))
	if err := h.svc.Approve(c.Param("id"), approverID); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "pelanggaran berhasil disetujui", nil)
}

func (h *PunishmentHandler) Reject(c *gin.Context) {
	uid, _ := c.Get("user_id")
	approverID := int64(uid.(float64))
	var req rejectReq
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.Reject(c.Param("id"), approverID, req.Reason); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "pelanggaran berhasil ditolak", nil)
}

func (h *PunishmentHandler) Create(c *gin.Context) {
	var input models.PunishmentLog
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Create(&input); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, "pelanggaran berhasil dicatat", input)
}

func (h *PunishmentHandler) Update(c *gin.Context) {
	var input models.PunishmentLog
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Update(c.Param("id"), &input); err != nil {
		response.NotFound(c, "catatan tidak ditemukan")
		return
	}
	response.OK(c, "catatan berhasil diperbarui", nil)
}

func (h *PunishmentHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "catatan berhasil dihapus", nil)
}
