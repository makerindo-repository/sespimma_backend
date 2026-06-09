package controllers

import (
	"net/http"

	"sespima_api/config"
	"sespima_api/dto"
	"sespima_api/models"
	

	"github.com/gin-gonic/gin"
)

type AkademikController struct{}

func NewAkademikController() *AkademikController {
	return &AkademikController{}
}

func buildAkademikTree(
	components []models.Akademik,
	nilaiMap map[uint64]models.PenilaianAkademik,
) []dto.AkademikNodeResponse {

	var result []dto.AkademikNodeResponse

	for _, comp := range components {

		node := dto.AkademikNodeResponse{
			ID:                  uint64(comp.ID),
			AkademikComponentID: uint64(comp.ID),
			Code:                comp.Code,
			Name:                comp.Name,
			Weight:              comp.Weight,
		}

		if nilaiMap != nil {
			if nilai, ok := nilaiMap[uint64(comp.ID)]; ok {
				node.Nilai = &nilai.Nilai
				node.Catatan = nilai.Catatan
			}
		}

		if len(comp.Children) > 0 {
			node.Children = buildAkademikTree(comp.Children, nilaiMap)
		}

		result = append(result, node)
	}

	return result
}

func (a *AkademikController) GetAll(c *gin.Context) {
	var components []models.Akademik

	if err := config.DB.
		Preload("Children.Children.Children.Children").
		Where("parent_id IS NULL").
		Find(&components).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := buildAkademikTree(components, nil)

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (a *AkademikController) GetBySerdikID(c *gin.Context) {
	serdikID := c.Param("serdikId")

	var components []models.Akademik

	if err := config.DB.
		Preload("Children.Children.Children.Children").
		Where("parent_id IS NULL").
		Find(&components).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	var penilaian []models.PenilaianAkademik

	config.DB.
		Where("serdik_id = ?", serdikID).
		Find(&penilaian)

	nilaiMap := make(map[uint64]models.PenilaianAkademik)

	for _, item := range penilaian {
		nilaiMap[item.AkademikComponentID] = item
	}

	result := buildAkademikTree(components, nilaiMap)

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (a *AkademikController) Create(c *gin.Context) {
	var req dto.PenilaianAkademikRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	data := models.PenilaianAkademik{
		SerdikID:            req.SerdikID,
		AkademikComponentID: req.AkademikComponentID,
		Nilai:               req.Nilai,
		Catatan:             req.Catatan,
	}

	if err := config.DB.Create(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"data":    data,
	})
}

func (a *AkademikController) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.PenilaianAkademikRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var data models.PenilaianAkademik

	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	data.Nilai = req.Nilai
	data.Catatan = req.Catatan

	if err := config.DB.Save(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    data,
	})
}

func (a *AkademikController) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.PenilaianAkademik{}, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}