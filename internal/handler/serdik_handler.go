package handler

import (
	"github.com/gin-gonic/gin"

	"sespima_api/config"
	"sespima_api/models"
	"sespima_api/pkg/response"
)

// SerdikHandler exposes the authenticated Serdik's own profile.
// Routes are role-guarded to "serdik" in the router.
type SerdikHandler struct{}

func NewSerdikHandler() *SerdikHandler { return &SerdikHandler{} }

// userIDFromCtx reads the user_id the auth middleware stored (as float64).
func userIDFromCtx(c *gin.Context) (uint, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	f, ok := v.(float64)
	if !ok {
		return 0, false
	}
	return uint(f), true
}

// GetMe returns the authenticated Serdik's full profile in the same shape
// the mobile app already consumes from /api/me and the login response.
func (h *SerdikHandler) GetMe(c *gin.Context) {
	userID, ok := userIDFromCtx(c)
	if !ok {
		response.Unauthorized(c, "sesi tidak valid")
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		response.NotFound(c, "user tidak ditemukan")
		return
	}

	data := buildLoginResponse(&user, "", "")
	delete(data, "access_token")
	delete(data, "refresh_token")
	response.OK(c, "success", data)
}

// updateSerdikReq holds the only Serdik fields a serdik may self-edit.
// Academic identity fields (NRP, NoSerdik, Pangkat, Pokjar, …) are immutable
// from the app and may only be changed by an administrator.
type updateSerdikReq struct {
	Email       *string `json:"email"`
	NoTelepon   *string `json:"no_telepon"`
	NoHandphone *string `json:"no_handphone"`
	Alamat      *string `json:"alamat"`
}

// UpdateMe updates the authenticated Serdik's editable contact fields.
func (h *SerdikHandler) UpdateMe(c *gin.Context) {
	userID, ok := userIDFromCtx(c)
	if !ok {
		response.Unauthorized(c, "sesi tidak valid")
		return
	}

	var req updateSerdikReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var serdik models.Serdik
	if err := config.DB.Where("user_id = ?", userID).First(&serdik).Error; err != nil {
		response.NotFound(c, "data serdik tidak ditemukan untuk akun ini")
		return
	}

	// Apply only the fields that were provided (non-nil).
	updates := map[string]interface{}{}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.NoTelepon != nil {
		updates["no_telepon"] = *req.NoTelepon
	}
	if req.NoHandphone != nil {
		updates["no_handphone"] = *req.NoHandphone
	}
	if req.Alamat != nil {
		updates["alamat"] = *req.Alamat
	}

	if len(updates) == 0 {
		response.BadRequest(c, "tidak ada perubahan yang dikirim")
		return
	}

	if err := config.DB.Model(&serdik).Updates(updates).Error; err != nil {
		response.InternalError(c, "gagal memperbarui profil")
		return
	}

	// Return the refreshed profile in the canonical shape.
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		response.OK(c, "profil berhasil diperbarui", nil)
		return
	}
	data := buildLoginResponse(&user, "", "")
	delete(data, "access_token")
	delete(data, "refresh_token")
	response.OK(c, "profil berhasil diperbarui", data)
}
