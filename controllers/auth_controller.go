package controllers

import (
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
    Email   string `json:"email"`
    NrpNip  string `json:"nrp_nip"`
    Password string `json:"password" binding:"required"`
}


func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    query := config.DB

    // Flexible lookup: email OR nrp_nip
    if req.Email != "" {
        query = query.Where("email = ?", req.Email)
    } else if req.NrpNip != "" {
        if len(req.NrpNip) == 8 {
            query = query.Where("nrp = ?", req.NrpNip)
        } else if len(req.NrpNip) > 8 {
            query = query.Where("nip = ?", req.NrpNip)
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nrp_nip format"})
            return
        }
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": "must provide email or nrp_nip"})
        return
    }

    if err := query.First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    if !user.IsActive {
        c.JSON(http.StatusForbidden, gin.H{"error": "account inactive"})
        return
    }

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    // Generate JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "nrp":     user.NRP,
        "nip":     user.NIP,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
        return
    }

    user.CurrentToken = &tokenString
    config.DB.Save(&user)

    c.JSON(http.StatusOK, gin.H{
        "message": "login successful",
        "token":   tokenString,
        "user": gin.H{
            "id":             user.ID,
            "email":          user.Email,
            "nrp":            user.NRP,
            "nip":            user.NIP,
            "role":           user.Role,
            "is_first_login": user.IsFirstLogin,
        },
    })
}


type GenerateTokenRequest struct {
    NrpNip string `json:"nrp_nip" binding:"required"`
}

func GenerateResetToken(c *gin.Context) {
    // Only operator should call this
    role, exists := c.Get("role")
    if !exists || role != "operator" {
        c.JSON(http.StatusForbidden, gin.H{"error": "only operator can generate reset token"})
        return
    }

    var req GenerateTokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := config.DB.Where("nrp = ? OR nip = ?", req.NrpNip, req.NrpNip).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    // Generate random 6-digit token
    // For simplicity, just generating a static one or pseudo-random
    token := "123456" // In production, use crypto/rand
    
    user.ResetToken = &token
    config.DB.Save(&user)

    c.JSON(http.StatusOK, gin.H{
        "message": "reset token generated successfully",
        "token": token,
    })
}

type ResetPasswordRequest struct {
    NrpNip   string `json:"nrp_nip" binding:"required"`
    Token    string `json:"token" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func ResetPassword(c *gin.Context) {
    var req ResetPasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := config.DB.Where("nrp = ? OR nip = ?", req.NrpNip, req.NrpNip).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    if user.ResetToken == nil || *user.ResetToken != req.Token {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
        return
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
        return
    }

    user.Password = string(hash)
    user.ResetToken = nil // Clear token after use
    config.DB.Save(&user)

    c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

type UpdatePasswordRequest struct {
    NewPassword string `json:"new_password" binding:"required"`
}

func UpdatePassword(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    var req UpdatePasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := config.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
        return
    }

    user.Password = string(hash)
    user.IsFirstLogin = false // Once they update it, they are no longer on first login
    config.DB.Save(&user)

    c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}
