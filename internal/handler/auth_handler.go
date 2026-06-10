package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"sespima_api/config"
	"sespima_api/internal/service"
	"sespima_api/models"
	"sespima_api/pkg/response"
)

type AuthHandler struct{ svc *service.AuthService }

func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

type loginReq struct {
	NrpNip   string `json:"nrp_nip" binding:"required"`
	Password string `json:"password" binding:"required"`
	FcmToken string `json:"fcm_token"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.Login(req.NrpNip, req.Password, req.FcmToken)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() == "akun tidak aktif" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"status": false, "message": err.Error()})
		return
	}

	data := buildLoginResponse(result.User, result.AccessToken, result.RefreshToken)
	response.OK(c, "login berhasil", data)
}

type generateTokenReq struct {
	NrpNip string `json:"nrp_nip" binding:"required"`
}

func (h *AuthHandler) GenerateResetToken(c *gin.Context) {
	var req generateTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	token, err := h.svc.GenerateResetToken(req.NrpNip)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "reset token berhasil dibuat", gin.H{"token": token})
}

type verifyTokenReq struct {
	NrpNip string `json:"nrp_nip" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

func (h *AuthHandler) VerifyResetToken(c *gin.Context) {
	var req verifyTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.VerifyResetToken(req.NrpNip, req.Token); err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.OK(c, "token valid", nil)
}

type resetPasswordReq struct {
	NrpNip   string `json:"nrp_nip" binding:"required"`
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req resetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.ResetPassword(req.NrpNip, req.Token, req.Password); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "password berhasil direset", nil)
}

type updatePasswordReq struct {
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	uid, _ := c.Get("user_id")
	userID := uint(uid.(float64))

	var req updatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdatePassword(userID, req.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "password berhasil diperbarui", nil)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	uid, _ := c.Get("user_id")
	userID := uint(uid.(float64))
	user, err := h.svc.GetUserByID(userID)
	if err != nil {
		response.NotFound(c, "user tidak ditemukan")
		return
	}
	data := buildLoginResponse(user, "", "")
	delete(data, "access_token")
	delete(data, "refresh_token")
	response.OK(c, "success", data)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	uid, _ := c.Get("user_id")
	userID := uint(uid.(float64))
	if err := h.svc.Logout(userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "logout berhasil", nil)
}

type refreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.svc.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.OK(c, "token diperbarui", gin.H{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}

// ---- login response builder ----

func safeStr(s *string) string {
	if s != nil {
		return *s
	}
	return "-"
}

var homeRoutes = map[string]string{
	"serdik":       "serdik_home",
	"gadik":        "gadik_home",
	"patun":        "patun_home",
	"korsis":       "korsis_home",
	"pimpinan":     "pimpinan_dashboard",
	"medis":        "medis_home",
	"operator":     "operator_home",
	"kabag_bindik": "kabag_bindik_home",
}

func buildLoginResponse(user *models.User, accessToken, refreshToken string) gin.H {
	nrp := "-"
	if user.NRP != nil {
		nrp = *user.NRP
	}

	homeRoute, ok := homeRoutes[user.Role]
	if !ok {
		homeRoute = "serdik_home"
	}

	data := gin.H{
		"user_id":             fmt.Sprintf("%d", user.ID),
		"name":                "-",
		"role_id":             user.Role,
		"pokjar":              "-",
		"nrp":                 nrp,
		"nosis":               "-",
		"pangkat":             "-",
		"angkatan":            "-",
		"agama":               "-",
		"jabatan":             "-",
		"no_serdik":           "-",
		"nik":                 "-",
		"jabatan_senat":       "-",
		"tempat_lahir":        "-",
		"no_handphone":        "-",
		"pendidikan_terakhir": "-",
		"alamat_lengkap":      "-",
		"email":               user.Email,
		"no_telepon":          "-",
		"kelompok":            "-",
		"diktuk_awal":         "-",
		"tahun_diktuk":        "-",
		"personel":            "-",
		"satker":              "-",
		"eselon":              "-",
		"golongan":            "-",
		"tanggal_lahir":       "1990-01-01",
		"jenis_kelamin":       "-",
		"is_first_login":      user.IsFirstLogin,
		"is_nak_approved":     user.IsNakApproved,
		"nilai_akademik":      0.0,
		"nilai_mental":        0.0,
		"nilai_jasmani":       0.0,
		"profile_photo":       user.ProfilePhoto,
		"access_token":        accessToken,
		"refresh_token":       refreshToken,
		"home_route":          homeRoute,
	}

	db := config.DB

	switch user.Role {
	case "serdik":
		var serdik models.Serdik
		if err := db.Where("user_id = ?", user.ID).Preload("Pokjar").First(&serdik).Error; err == nil {
			data["name"] = safeStr(serdik.NamaLengkap)
			data["no_serdik"] = safeStr(serdik.NoSerdik)
			data["nosis"] = safeStr(serdik.Nosis)
			data["nik"] = safeStr(serdik.NIK)
			data["nrp"] = safeStr(serdik.NRP)
			data["pangkat"] = safeStr(serdik.Pangkat)
			data["jabatan"] = safeStr(serdik.Jabatan)
			data["jabatan_senat"] = safeStr(serdik.JabatanSenat)
			data["agama"] = safeStr(serdik.Agama)
			data["tempat_lahir"] = safeStr(serdik.TempatLahir)
			data["no_handphone"] = safeStr(serdik.NoHandphone)
			data["no_telepon"] = safeStr(serdik.NoTelepon)
			data["alamat_lengkap"] = safeStr(serdik.Alamat)
			data["email"] = safeStr(serdik.Email)
			data["pendidikan_terakhir"] = safeStr(serdik.PendidikanTerakhir)
			data["diktuk_awal"] = safeStr(serdik.DiktukAwal)
			data["satker"] = safeStr(serdik.Satker)
			data["jenis_kelamin"] = serdik.JenisKelamin
			data["personel"] = fmt.Sprintf("%t", serdik.IsPersonel)
			data["angkatan"] = safeStr(serdik.Angkatan)
			if serdik.KelompokKelas != nil {
				data["kelompok"] = *serdik.KelompokKelas
				data["pokjar"] = *serdik.KelompokKelas
			}
			if serdik.TahunDiktuk != nil {
				data["tahun_diktuk"] = fmt.Sprintf("%d", *serdik.TahunDiktuk)
			}
			if serdik.TanggalLahir != nil {
				data["tanggal_lahir"] = serdik.TanggalLahir.Format("2006-01-02")
			}
			if serdik.Pokjar != nil {
				data["pokjar"] = serdik.Pokjar.Name
			}
		}
	case "gadik":
		var gadik models.Gadik
		if err := db.Where("user_id = ?", user.ID).First(&gadik).Error; err == nil {
			data["name"] = gadik.Nama
			data["pangkat"] = gadik.Pangkat
			data["jabatan"] = safeStr(gadik.JabatanStruktural)
			data["agama"] = safeStr(gadik.Agama)
			data["eselon"] = safeStr(gadik.Eselon)
			data["golongan"] = safeStr(gadik.Golongan)
			if gadik.NrpNip != nil {
				data["nrp"] = *gadik.NrpNip
			}
		}
	case "patun":
		var patun models.Patun
		if err := db.Where("user_id = ?", user.ID).First(&patun).Error; err == nil {
			data["name"] = patun.Nama
			data["pangkat"] = safeStr(patun.Pangkat)
			data["jabatan"] = safeStr(patun.JabatanStruktural)
			data["jabatan_senat"] = safeStr(patun.PeranPengasuhan)
			data["nrp"] = patun.NrpNip
			if patun.Pokjar != nil {
				data["pokjar"] = *patun.Pokjar
			}
		}
	case "korsis":
		var korsis models.Korsis
		if err := db.Where("user_id = ?", user.ID).First(&korsis).Error; err == nil {
			data["name"] = korsis.Nama
			data["pangkat"] = korsis.Pangkat
			data["jabatan"] = safeStr(korsis.JabatanStruktural)
			data["jabatan_senat"] = safeStr(korsis.PeranPengasuhan)
			if korsis.NrpNip != nil {
				data["nrp"] = *korsis.NrpNip
			}
		}
	case "pimpinan":
		var pimpinan models.Pimpinan
		if err := db.Where("user_id = ?", user.ID).First(&pimpinan).Error; err == nil {
			data["name"] = pimpinan.Nama
			data["pangkat"] = pimpinan.Pangkat
			data["jabatan"] = pimpinan.JabatanStruktural
			data["jabatan_senat"] = pimpinan.PeranPengasuhan
			data["nrp"] = pimpinan.NrpNip
		}
	case "medis":
		data["name"] = "Dokter SESPIMMA"
		data["jabatan"] = "Dokter Poliklinik SESPIMMA"
		if user.NIP != nil {
			data["nrp"] = *user.NIP
		}
	case "operator":
		data["name"] = "Operator SESPIMMA"
		data["jabatan"] = "Operator"
	case "kabag_bindik":
		data["name"] = "Kabag Bindik SESPIMMA"
		data["jabatan"] = "Kepala Bagian Bimbingan Pendidikan"
	}

	return data
}
