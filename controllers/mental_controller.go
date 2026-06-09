package controllers

import (
	"strconv"

	"sespima_api/config"
	"sespima_api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MentalNode struct {
	ID            uint64       `json:"id"`
	Code          string       `json:"code"`
	Name          string       `json:"name"`
	Weight        float64      `json:"weight"`
	IndicatorType string       `json:"indicator_type"`
	Point         *float64     `json:"point,omitempty"`
	Description   string       `json:"description"`
	Nilai         *float64     `json:"nilai,omitempty"`
	Catatan       *string      `json:"catatan,omitempty"`
	Children      []MentalNode `json:"children,omitempty"`
}

type CreatePenilaianMentalRequest struct {
	SerdikID          uint64  `json:"serdik_id"`
	MentalComponentID uint64  `json:"mental_component_id"`
	Nilai             float64 `json:"nilai"`
	Catatan           *string `json:"catatan"`
}

type UpdatePenilaianMentalRequest struct {
	Nilai   float64 `json:"nilai"`
	Catatan *string `json:"catatan"`
}

type MentalController struct{}

func NewMentalController() *MentalController {
	return &MentalController{}
}

// =========================
// GET ALL
// =========================

func (mc *MentalController) GetAll(c *fiber.Ctx) error {
	var components []models.Mental

	if err := config.DB.
		Preload("Children.Children.Children").
		Order("id ASC").
		Find(&components).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	tree := buildMentalTree(components, nil)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    tree,
	})
}

// =========================
// GET BY SERDIK ID
// =========================

func (mc *MentalController) GetBySerdikID(c *fiber.Ctx) error {
	serdikID, err := strconv.ParseUint(c.Params("serdikId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid serdik id",
		})
	}

	var components []models.Mental
	if err := config.DB.Find(&components).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	var penilaian []models.PenilaianMental
	if err := config.DB.
		Where("serdik_id = ?", serdikID).
		Find(&penilaian).Error; err != nil {

		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	nilaiMap := make(map[uint64]models.PenilaianMental)
	for _, p := range penilaian {
		nilaiMap[p.MentalComponentID] = p
	}

	tree := buildMentalTreeWithScore(components, nilaiMap, nil)

	return c.JSON(fiber.Map{
		"success": true,
		"data":    tree,
	})
}

// =========================
// CREATE
// =========================

func (mc *MentalController) Create(c *fiber.Ctx) error {
	var req CreatePenilaianMentalRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data := models.PenilaianMental{
		SerdikID:          req.SerdikID,
		MentalComponentID: req.MentalComponentID,
		Nilai:             req.Nilai,
		Catatan:           req.Catatan,
	}

	if err := config.DB.Create(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "penilaian mental created",
		"data":    data,
	})
}

// =========================
// UPDATE
// =========================

func (mc *MentalController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid id",
		})
	}

	var req UpdatePenilaianMentalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	var data models.PenilaianMental

	if err := config.DB.First(&data, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "data not found",
		})
	}

	data.Nilai = req.Nilai
	data.Catatan = req.Catatan

	if err := config.DB.Save(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "penilaian mental updated",
		"data":    data,
	})
}

// =========================
// DELETE
// =========================

func (mc *MentalController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid id",
		})
	}

	if err := config.DB.Delete(&models.PenilaianMental{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "penilaian mental deleted",
	})
}

// =========================
// TREE BUILDER
// =========================

func buildMentalTree(
	components []models.Mental,
	parentID *int64,
) []MentalNode {

	var result []MentalNode

	for _, component := range components {

		if !sameParent(component.ParentID, parentID) {
			continue
		}

		node := MentalNode{
			ID:            uint64(component.ID),
			Code:          component.Code,
			Name:          component.Name,
			Weight:        component.Weight,
			IndicatorType: string(component.IndicatorType),
			Point:         component.Point,
			Description:   component.Description,
			Children:      buildMentalTree(components, &component.ID),
		}

		result = append(result, node)
	}

	return result
}

func buildMentalTreeWithScore(
	components []models.Mental,
	nilaiMap map[uint64]models.PenilaianMental,
	parentID *int64,
) []MentalNode {

	var result []MentalNode

	for _, component := range components {

		if !sameParent(component.ParentID, parentID) {
			continue
		}

		node := MentalNode{
			ID:            uint64(component.ID),
			Code:          component.Code,
			Name:          component.Name,
			Weight:        component.Weight,
			IndicatorType: string(component.IndicatorType),
			Point:         component.Point,
			Description:   component.Description,
		}

		if nilai, ok := nilaiMap[uint64(component.ID)]; ok {
			node.Nilai = &nilai.Nilai
			node.Catatan = nilai.Catatan
		}

		node.Children = buildMentalTreeWithScore(
			components,
			nilaiMap,
			&component.ID,
		)

		result = append(result, node)
	}

	return result
}

func sameParent(a, b *int64) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return *a == *b
}

func isNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}