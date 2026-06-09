package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sespima_api/config"
	"sespima_api/models"
)

type UserLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// GET /users
func GetAllUsers(c *gin.Context) {
	var users []models.User

	if err := config.DB.
		Preload("Pimpinan").
		Preload("Patun").
		Preload("Korsis").
		Preload("Serdik").
		Preload("Gadik").
		Find(&users).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed get users",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    users,
	})
}

// GET /users/:id/location
func GetUserLocation(c *gin.Context) {
	id := c.Param("id")

	type Result struct {
		ID        uint     `json:"id"`
		Email     string   `json:"email"`
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
	}

	var result Result

	err := config.DB.Raw(`
		SELECT
			id,
			email,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude
		FROM users
		WHERE id = ?
	`, id).Scan(&result).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed get location",
			"error":   err.Error(),
		})
		return
	}

	if result.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    result,
	})
}

// PUT /users/:id/location
func UpdateUserLocation(c *gin.Context) {
	id := c.Param("id")

	var req UserLocationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	result := config.DB.Exec(`
		UPDATE users
		SET location = ST_SetSRID(
			ST_MakePoint(?, ?),
			4326
		)::geography
		WHERE id = ?
	`, req.Longitude, req.Latitude, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed update location",
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "location updated",
	})
}

// OPTIONAL - update by authenticated user
func UpdateMyLocation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	var req UserLocationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	result := config.DB.Exec(`
		UPDATE users
		SET location = ST_SetSRID(
			ST_MakePoint(?, ?),
			4326
		)::geography
		WHERE id = ?
	`, req.Longitude, req.Latitude, userID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed update location",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "location updated",
	})
}

// GET /users/locations
func GetAllUserLocations(c *gin.Context) {
	type UserLocation struct {
		ID        uint     `json:"id"`
		Email     string   `json:"email"`
		Role      string   `json:"role"`
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
	}

	var locations []UserLocation

	err := config.DB.Raw(`
		SELECT
			id,
			email,
			role,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude
		FROM users
		WHERE location IS NOT NULL
		ORDER BY id
	`).Scan(&locations).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed get locations",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    locations,
	})
}