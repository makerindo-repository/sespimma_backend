package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "sespima_api/models"
)

// =========================
// PunishmentLog Controller
// =========================

type PunishmentLogController struct {
    DB *gorm.DB
}

// GetAll - fetch all PunishmentLogs
func (c *PunishmentLogController) GetAll(ctx *gin.Context) {
    var logs []models.PunishmentLog
    if err := c.DB.Preload("PunishmentItem").Find(&logs).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, logs)
}

// GetByUserID - fetch PunishmentLogs by UserID
func (c *PunishmentLogController) GetByUserID(ctx *gin.Context) {
    userID := ctx.Param("user_id")
    var logs []models.PunishmentLog
    if err := c.DB.Where("user_id = ?", userID).Preload("PunishmentItem").Find(&logs).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, logs)
}

// Create - insert new PunishmentLog
func (c *PunishmentLogController) Create(ctx *gin.Context) {
    var input models.PunishmentLog
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

// Update - update existing PunishmentLog
func (c *PunishmentLogController) Update(ctx *gin.Context) {
    id := ctx.Param("id")
    var log models.PunishmentLog
    if err := c.DB.First(&log, id).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "PunishmentLog not found"})
        return
    }

    var input models.PunishmentLog
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.DB.Model(&log).Updates(input).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, log)
}

// Delete - remove PunishmentLog
func (c *PunishmentLogController) Delete(ctx *gin.Context) {
    id := ctx.Param("id")
    if err := c.DB.Delete(&models.PunishmentLog{}, id).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "PunishmentLog deleted"})
}
