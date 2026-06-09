package controllers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "sespima_api/models"
)

type KegiatanController struct {
    DB *gorm.DB
}

// =========================
// Get All Kegiatan
// =========================
func (c *KegiatanController) GetAll(ctx *gin.Context) {
    var kegiatan []models.Kegiatan
    if err := c.DB.Preload("Lokasi").Find(&kegiatan).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, kegiatan)
}

// =========================
// Get Kegiatan by IsRutin = true
// =========================
func (c *KegiatanController) GetRutin(ctx *gin.Context) {
    var kegiatan []models.Kegiatan
    if err := c.DB.Preload("Lokasi").Where("is_rutin = ?", true).Find(&kegiatan).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, kegiatan)
}

// =========================
// Get Kegiatan by TanggalMulai = today
// =========================
func (c *KegiatanController) GetToday(ctx *gin.Context) {
    today := time.Now().Format("2006-01-02")
    var kegiatan []models.Kegiatan
    if err := c.DB.Preload("Lokasi").Where("tanggal_mulai = ?", today).Find(&kegiatan).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, kegiatan)
}

// Create Kegiatan
func (c *KegiatanController) Create(ctx *gin.Context) {
    var input models.Kegiatan
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.DB.Create(&input).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, input)
}

// Update Kegiatan
func (c *KegiatanController) Update(ctx *gin.Context) {
    id := ctx.Param("id")
    var kegiatan models.Kegiatan
    if err := c.DB.First(&kegiatan, id).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "Kegiatan not found"})
        return
    }

    var input models.Kegiatan
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    kegiatan.Nama = input.Nama
    kegiatan.NamaLokasi = input.NamaLokasi
    kegiatan.TanggalMulai = input.TanggalMulai
    kegiatan.WaktuMulai = input.WaktuMulai
    kegiatan.WaktuSelesai = input.WaktuSelesai
    kegiatan.BatasWaktuPenugasan = input.BatasWaktuPenugasan
    kegiatan.IsRutin = input.IsRutin
    kegiatan.IsPelatihan = input.IsPelatihan
    kegiatan.QRCode = input.QRCode
    kegiatan.Radius = input.Radius
    kegiatan.Location = input.Location // bisa POINT atau POLYGON

    if err := c.DB.Save(&kegiatan).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, kegiatan)
}

// =========================
// Delete Kegiatan
// =========================
func (c *KegiatanController) Delete(ctx *gin.Context) {
    id := ctx.Param("id")
    if err := c.DB.Delete(&models.Kegiatan{}, id).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "Kegiatan deleted"})
}
