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
