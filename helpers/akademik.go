package helpers

import (
	"sespima_api/dto"
	"sespima_api/models"
)

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