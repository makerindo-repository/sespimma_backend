package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sespima_api/internal/service"
	"sespima_api/models"
	"sespima_api/pkg/response"
)

type AssignmentHandler struct{ svc *service.AssignmentService }

func NewAssignmentHandler(svc *service.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{svc: svc}
}

func (h *AssignmentHandler) GetAll(c *gin.Context) {
	var pokjarID *int
	if v := c.Query("pokjar_id"); v != "" {
		id, _ := strconv.Atoi(v)
		pokjarID = &id
	}
	list, err := h.svc.GetAll(pokjarID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *AssignmentHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	a, err := h.svc.GetByID(id)
	if err != nil {
		response.NotFound(c, "tugas tidak ditemukan")
		return
	}
	response.OK(c, "success", a)
}

func (h *AssignmentHandler) Create(c *gin.Context) {
	var input models.Assignment
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	uid, _ := c.Get("user_id")
	input.CreatedBy = int64(uid.(float64))
	if err := h.svc.Create(&input); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, "tugas berhasil dibuat", input)
}

func (h *AssignmentHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var input models.Assignment
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Update(id, &input); err != nil {
		response.NotFound(c, "tugas tidak ditemukan")
		return
	}
	response.OK(c, "tugas berhasil diperbarui", nil)
}

func (h *AssignmentHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "tugas berhasil dihapus", nil)
}

func (h *AssignmentHandler) GetSubmissions(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	list, err := h.svc.GetSubmissions(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", list)
}

func (h *AssignmentHandler) Submit(c *gin.Context) {
	var input models.Submission
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Submit(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "tugas berhasil dikumpulkan", input)
}

func (h *AssignmentHandler) GradeSubmission(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("submissionId"), 10, 64)
	var input models.Submission
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	input.ID = id
	input.IsGraded = true
	if err := h.svc.GradeSubmission(&input); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "nilai berhasil disimpan", nil)
}
