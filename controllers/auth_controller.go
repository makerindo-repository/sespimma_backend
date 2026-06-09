package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"sespima_api/config"
	"sespima_api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	NrpNip   string `json:"nrp_nip" binding:"required"`
	Password string `json:"password" binding:"required"`
	FcmToken string `json:"fcm_token"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	var user models.User
	query := config.DB

	if len(req.NrpNip) > 10 {
		query = query.Where("nip = ?", req.NrpNip)
	} else {
		query = query.Where("nrp = ?", req.NrpNip)
	}

	if err := query.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "NRP/NIP atau password salah"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"status": false, "message": "Akun tidak aktif"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "NRP/NIP atau password salah"})
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"nrp":     user.NRP,
		"nip":     user.NIP,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "gagal membuat token"})
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "gagal membuat refresh token"})
		return
	}

	user.CurrentToken = &accessTokenString
	config.DB.Save(&user)

	data := buildLoginResponseData(user, accessTokenString, refreshTokenString)

	c.JSON(http.StatusOK, gin.H{"status": true, "data": data})
}

func buildLoginResponseData(user models.User, accessToken string, refreshToken string) gin.H {
	nrpStr := "-"
	if user.NRP != nil {
		nrpStr = *user.NRP
	}

	emailStr := "-"
	if user.Email != "" {
		emailStr = user.Email
	}

	data := gin.H{
		"user_id":              fmt.Sprintf("%d", user.ID),
		"name":                 "-",
		"role_id":              user.Role,
		"pokjar":               "-",
		"nrp":                  nrpStr,
		"nosis":                "-",
		"pangkat":              "-",
		"angkatan":             "-",
		"agama":                "-",
		"jabatan":              "-",
		"no_serdik":            "-",
		"nik":                  "-",
		"jabatan_senat":        "-",
		"tempat_lahir":         "-",
		"no_handphone":         "-",
		"pendidikan_terakhir":  "-",
		"alamat_lengkap":       "-",
		"email":                emailStr,
		"no_telepon":           "-",
		"kelompok":             "-",
		"diktuk_awal":          "-",
		"tahun_diktuk":         "-",
		"personel":             "-",
		"satker":               "-",
		"tanggal_lahir":        "1990-01-01",
		"eselon":               "-",
		"golongan":             "-",
		"is_nak_approved":      false,
		"is_first_login":       user.IsFirstLogin,
		"jenis_kelamin":        "-",
		"nilai_akademik":       0.0,
		"nilai_mental":         0.0,
		"nilai_jasmani":        0.0,
		"access_token":         accessToken,
		"refresh_token":        refreshToken,
	}

	switch user.Role {
	case "serdik":
		var serdik models.Serdik
		if err := config.DB.Where("user_id = ?", user.ID).Preload("Pokjar").First(&serdik).Error; err == nil {
			data["name"] = safeStr(serdik.NamaLengkap)
			data["no_serdik"] = safeStr(serdik.NoSerdik)
			data["nik"] = safeStr(serdik.NIK)
			data["nrp"] = safeStr(serdik.NRP)
			data["pangkat"] = safeStr(serdik.Pangkat)
			data["jabatan"] = safeStr(serdik.Jabatan)
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
			if serdik.KelompokKelas != nil {
				data["kelompok"] = *serdik.KelompokKelas
				data["pokjar"] = *serdik.KelompokKelas
			}
			if serdik.TahunDiktuk != nil {
				data["tahun_diktuk"] = fmt.Sprintf("%d", *serdik.TahunDiktuk)
				data["angkatan"] = fmt.Sprintf("%d", *serdik.TahunDiktuk)
			}
			if serdik.TanggalLahir != nil {
				data["tanggal_lahir"] = serdik.TanggalLahir.Format("2006-01-02")
			}
		}
	case "gadik":
		var gadik models.Gadik
		if err := config.DB.Where("user_id = ?", user.ID).First(&gadik).Error; err == nil {
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
		if err := config.DB.Where("user_id = ?", user.ID).First(&patun).Error; err == nil {
			data["name"] = patun.Nama
			data["pangkat"] = safeStrVal(patun.Pangkat)
			data["jabatan"] = safeStr(patun.JabatanStruktural)
			data["jabatan_senat"] = safeStr(patun.PeranPengasuhan)
			data["nrp"] = patun.NrpNip
			if patun.Pokjar != nil {
				data["pokjar"] = *patun.Pokjar
			}
		}
	case "korsis":
		var korsis models.Korsis
		if err := config.DB.Where("user_id = ?", user.ID).First(&korsis).Error; err == nil {
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
		if err := config.DB.Where("user_id = ?", user.ID).First(&pimpinan).Error; err == nil {
			data["name"] = pimpinan.Nama
			data["pangkat"] = pimpinan.Pangkat
			data["jabatan"] = pimpinan.JabatanStruktural
			data["jabatan_senat"] = pimpinan.PeranPengasuhan
			data["nrp"] = pimpinan.NrpNip
		}
	case "operator":
		data["name"] = "Operator Tim SESPIMMA"
		data["jabatan"] = "Operator"
	case "medis":
		data["name"] = "dr. Fenny Lestari"
		data["pangkat"] = "Penata"
		data["jabatan"] = "Dokter Poliklinik SESPIMMA"
		if user.NIP != nil {
			data["nrp"] = *user.NIP
		}
	}

	return data
}

func safeStr(s *string) string {
	if s != nil {
		return *s
	}
	return "-"
}

func safeStrVal(s *string) string {
	if s != nil {
		return *s
	}
	return "-"
}

func generateResetToken() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

type GenerateTokenRequest struct {
	NrpNip string `json:"nrp_nip" binding:"required"`
}

func GenerateResetToken(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "operator" {
		c.JSON(http.StatusForbidden, gin.H{"status": false, "message": "hanya operator yang dapat membuat reset token"})
		return
	}

	var req GenerateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("nrp = ? OR nip = ?", req.NrpNip, req.NrpNip).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "user tidak ditemukan"})
		return
	}

	token := generateResetToken()
	user.ResetToken = &token
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "reset token berhasil dibuat",
		"token":   token,
	})
}

type VerifyResetTokenRequest struct {
	NrpNip string `json:"nrp_nip" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

func VerifyResetToken(c *gin.Context) {
	var req VerifyResetTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("nrp = ? OR nip = ?", req.NrpNip, req.NrpNip).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "NRP tidak ditemukan"})
		return
	}

	if user.ResetToken == nil || *user.ResetToken != req.Token {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "token tidak valid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "token valid"})
}

type ResetPasswordRequest struct {
	NrpNip   string `json:"nrp_nip" binding:"required"`
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("nrp = ? OR nip = ?", req.NrpNip, req.NrpNip).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "user tidak ditemukan"})
		return
	}

	if user.ResetToken == nil || *user.ResetToken != req.Token {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "token tidak valid"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "gagal mengenkripsi password"})
		return
	}

	user.Password = string(hash)
	user.ResetToken = nil
	user.IsFirstLogin = false
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "password berhasil direset"})
}

type UpdatePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required"`
}

func UpdatePassword(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "unauthorized"})
		return
	}

	userID := uint(userIDRaw.(float64))

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "user tidak ditemukan"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "gagal mengenkripsi password"})
		return
	}

	user.Password = string(hash)
	user.IsFirstLogin = false
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "password berhasil diperbarui"})
}
