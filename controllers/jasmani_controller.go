package controllers

import (
	"strconv"

	"sespima_api/config"
	"sespima_api/dto"
	"sespima_api/models"

	"github.com/gofiber/fiber/v2"
)

type JasmaniController struct{}

func NewJasmaniController() *JasmaniController {
	return &JasmaniController{}
}

func buildJasmaniTree(
	component models.Jasmani,
	nilaiMap map[uint64]models.PenilaianJasmani,
) dto.JasmaniResponse {

	var nilai *float64
	var catatan *string

	if pen, ok := nilaiMap[uint64(component.ID)]; ok {
		nilai = &pen.Nilai
		catatan = pen.Catatan
	}

	children := make([]dto.JasmaniResponse, 0)

	for _, child := range component.Children {
		children = append(children, buildJasmaniTree(child, nilaiMap))
	}

	return dto.JasmaniResponse{
		ID:       uint64(component.ID),
		Code:     component.Code,
		Name:     component.Name,
		Weight:   component.Weight,
		Nilai:    nilai,
		Catatan:  catatan,
		Children: children,
	}
}

func (ctl *JasmaniController) GetAll(c *fiber.Ctx) error {
	var data []models.PenilaianJasmani

	if err := config.DB.
		Preload("Serdik").
		Preload("JasmaniComponent").
		Find(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func (ctl *JasmaniController) GetBySerdikID(c *fiber.Ctx) error {
	serdikID, err := strconv.ParseUint(c.Params("serdikId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid serdik id",
		})
	}

	var components []models.Jasmani

	if err := config.DB.
		Where("parent_id IS NULL").
		Preload("Children.Children.Children.Children").
		Find(&components).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	var penilaian []models.PenilaianJasmani

	if err := config.DB.
		Where("serdik_id = ?", serdikID).
		Find(&penilaian).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	nilaiMap := make(map[uint64]models.PenilaianJasmani)

	for _, p := range penilaian {
		nilaiMap[p.JasmaniComponentID] = p
	}

	result := make([]dto.JasmaniResponse, 0)

	for _, component := range components {
		result = append(result, buildJasmaniTree(component, nilaiMap))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctl *JasmaniController) Create(c *fiber.Ctx) error {
	var payload models.PenilaianJasmani

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := config.DB.Create(&payload).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    payload,
	})
}

func (ctl *JasmaniController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid id",
		})
	}

	var data models.PenilaianJasmani

	if err := config.DB.First(&data, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "data not found",
		})
	}

	var payload models.PenilaianJasmani

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data.Nilai = payload.Nilai
	data.Catatan = payload.Catatan
	data.SerdikID = payload.SerdikID
	data.JasmaniComponentID = payload.JasmaniComponentID

	if err := config.DB.Save(&data).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func (ctl *JasmaniController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid id",
		})
	}

	if err := config.DB.Delete(&models.PenilaianJasmani{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "deleted successfully",
	})
}
