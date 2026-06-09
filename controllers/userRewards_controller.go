package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "sespima_api/models"
)

// =========================
// UserReward Controller
// =========================

type UserRewardController struct {
    DB *gorm.DB
}

// GetAll - fetch all UserRewards
func (c *UserRewardController) GetAll(ctx *gin.Context) {
    var rewards []models.UserReward
    if err := c.DB.Preload("RewardItem").Find(&rewards).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, rewards)
}

// GetPendingApproval - fetch UserRewards where ApprovedBy is NULL
func (c *UserRewardController) GetPendingApproval(ctx *gin.Context) {
    var rewards []models.UserReward
    if err := c.DB.Where("approved_by IS NULL").Preload("RewardItem").Find(&rewards).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, rewards)
}

// Create - insert new UserReward
func (c *UserRewardController) Create(ctx *gin.Context) {
    var input models.UserReward
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

// Update - update existing UserReward
func (c *UserRewardController) Update(ctx *gin.Context) {
    id := ctx.Param("id")
    var reward models.UserReward
    if err := c.DB.First(&reward, id).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "UserReward not found"})
        return
    }

    var input models.UserReward
    if err := ctx.ShouldBindJSON(&input); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.DB.Model(&reward).Updates(input).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, reward)
}

// Delete - remove UserReward
func (c *UserRewardController) Delete(ctx *gin.Context) {
    id := ctx.Param("id")
    if err := c.DB.Delete(&models.UserReward{}, id).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "UserReward deleted"})
}

// Approve - set ApprovedBy for a UserReward
func (c *UserRewardController) Approve(ctx *gin.Context) {
    id := ctx.Param("id")
    userIDStr := ctx.Query("user_id") // who approves
    userID, err := strconv.ParseInt(userIDStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
        return
    }

    var reward models.UserReward
    if err := c.DB.First(&reward, id).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "UserReward not found"})
        return
    }

    reward.ApprovedBy = &userID
    reward.Status = "approved"

    if err := c.DB.Save(&reward).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, reward)
}
