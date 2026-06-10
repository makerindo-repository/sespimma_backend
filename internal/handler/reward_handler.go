package handler

import (
	"github.com/gin-gonic/gin"

	"sespima_api/internal/service"
	"sespima_api/models"
	"sespima_api/pkg/response"
)

type RewardHandler struct{ svc *service.RewardService }

func NewRewardHandler(svc *service.RewardService) *RewardHandler { return &RewardHandler{svc: svc} }

func (h *RewardHandler) GetAll(c *gin.Context) {
	list, err := h.svc.GetAll()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *RewardHandler) GetPending(c *gin.Context) {
	list, err := h.svc.GetPending()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *RewardHandler) Create(c *gin.Context) {
	var input models.UserReward
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Create(&input); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, "reward berhasil dicatat", input)
}

func (h *RewardHandler) Update(c *gin.Context) {
	var input models.UserReward
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Update(c.Param("id"), &input); err != nil {
		response.NotFound(c, "reward tidak ditemukan")
		return
	}
	response.OK(c, "reward berhasil diperbarui", nil)
}

func (h *RewardHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "reward berhasil dihapus", nil)
}

func (h *RewardHandler) Approve(c *gin.Context) {
	uid, _ := c.Get("user_id")
	approverID := int64(uid.(float64))
	if err := h.svc.Approve(c.Param("id"), approverID); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "reward berhasil disetujui", nil)
}

type rejectReq struct {
	Reason string `json:"reason"`
}

func (h *RewardHandler) Reject(c *gin.Context) {
	uid, _ := c.Get("user_id")
	approverID := int64(uid.(float64))
	var req rejectReq
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.Reject(c.Param("id"), approverID, req.Reason); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "reward berhasil ditolak", nil)
}
