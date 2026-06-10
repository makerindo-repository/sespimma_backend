package handler

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"sespima_api/internal/service"
	"sespima_api/internal/ws"
	"sespima_api/models"
	"sespima_api/pkg/response"
)

type UserHandler struct {
	svc *service.UserService
	hub *ws.Hub
}

func NewUserHandler(svc *service.UserService, hub *ws.Hub) *UserHandler {
	return &UserHandler{svc: svc, hub: hub}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.svc.GetAllWithProfiles()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", users)
}

func (h *UserHandler) GetLocation(c *gin.Context) {
	result, err := h.svc.GetUserLocation(c.Param("id"))
	if err != nil {
		response.NotFound(c, "user tidak ditemukan")
		return
	}
	response.OK(c, "success", result)
}

func (h *UserHandler) GetAllLocations(c *gin.Context) {
	results, err := h.svc.GetAllUserLocations()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "success", results)
}

type locationReq struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Accuracy  *float64 `json:"accuracy"`
}

func (h *UserHandler) UpdateLocation(c *gin.Context) {
	var req locationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateLocation(c.Param("id"), req.Latitude, req.Longitude); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "location updated", nil)
}

func (h *UserHandler) UpdateMyLocation(c *gin.Context) {
	uid, _ := c.Get("user_id")
	var req locationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateLocation(uid, req.Latitude, req.Longitude); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	userIDFloat, _ := uid.(float64)
	userID := uint(userIDFloat)

	log := &models.UserLocationLog{
		UserID:     int64(userIDFloat),
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Accuracy:   req.Accuracy,
		RecordedAt: time.Now(),
	}
	_ = h.svc.LogLocation(log)

	// Broadcast to live WS viewers
	user, err := h.svc.GetByID(userID)
	if err == nil {
		accuracy := 0.0
		if req.Accuracy != nil {
			accuracy = *req.Accuracy
		}
		role, _ := c.Get("role")
		h.hub.Broadcast <- ws.LocationUpdate{
			Type:      "location_update",
			UserID:    userID,
			Role:      role.(string),
			Name:      user.Email,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
			Accuracy:  accuracy,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	response.OK(c, "location updated", nil)
}

var allowedPhotoExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}

func (h *UserHandler) UploadProfilePhoto(c *gin.Context) {
	uid, _ := c.Get("user_id")
	userID := uint(uid.(float64))

	file, err := c.FormFile("photo")
	if err != nil {
		response.BadRequest(c, "file foto diperlukan")
		return
	}
	if file.Size > 5*1024*1024 {
		response.BadRequest(c, "ukuran file maksimal 5MB")
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedPhotoExts[ext] {
		response.BadRequest(c, "format file harus jpg, jpeg, png, atau webp")
		return
	}

	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixNano(), ext)
	savePath := filepath.Join("uploads", "profile", filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.InternalError(c, "gagal menyimpan foto")
		return
	}

	scheme := c.Request.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
	}
	photoURL := fmt.Sprintf("%s://%s/uploads/profile/%s", scheme, c.Request.Host, filename)

	if err := h.svc.UpdateProfilePhoto(userID, photoURL); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "foto profil berhasil diperbarui", gin.H{"photo_url": photoURL})
}

func (h *UserHandler) DeleteProfilePhoto(c *gin.Context) {
	uid, _ := c.Get("user_id")
	userID := uint(uid.(float64))
	if err := h.svc.ClearProfilePhoto(userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "foto profil berhasil dihapus", nil)
}
