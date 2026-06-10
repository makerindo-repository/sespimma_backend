package service

import (
	"sespima_api/dto"
	"sespima_api/internal/repository"
	"sespima_api/models"
)

type AssessmentService struct {
	akademik  repository.AkademikRepository
	jasmani   repository.JasmaniRepository
	mental    repository.MentalRepository
	sociometry repository.SociometryRepository
	health    repository.HealthRepository
}

func NewAssessmentService(
	akademik repository.AkademikRepository,
	jasmani repository.JasmaniRepository,
	mental repository.MentalRepository,
	sociometry repository.SociometryRepository,
	health repository.HealthRepository,
) *AssessmentService {
	return &AssessmentService{
		akademik:   akademik,
		jasmani:    jasmani,
		mental:     mental,
		sociometry: sociometry,
		health:     health,
	}
}

// ---- Akademik ----

func (s *AssessmentService) GetAkademikTree() ([]dto.AkademikNodeResponse, error) {
	components, err := s.akademik.GetAllComponents()
	if err != nil {
		return nil, err
	}
	return buildAkademikTree(components, nil), nil
}

func (s *AssessmentService) GetAkademikBySerdik(serdikID string) ([]dto.AkademikNodeResponse, error) {
	components, err := s.akademik.GetAllComponents()
	if err != nil {
		return nil, err
	}
	penilaian, err := s.akademik.GetPenilaianBySerdik(serdikID)
	if err != nil {
		return nil, err
	}
	nilaiMap := make(map[uint64]models.PenilaianAkademik)
	for _, p := range penilaian {
		nilaiMap[p.AkademikComponentID] = p
	}
	return buildAkademikTree(components, nilaiMap), nil
}

func buildAkademikTree(components []models.Akademik, nilaiMap map[uint64]models.PenilaianAkademik) []dto.AkademikNodeResponse {
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

func (s *AssessmentService) CreateAkademikPenilaian(p *models.PenilaianAkademik) error {
	return s.akademik.CreatePenilaian(p)
}

func (s *AssessmentService) UpdateAkademikPenilaian(id string, nilai float64, catatan *string) error {
	return s.akademik.UpdatePenilaian(id, nilai, catatan)
}

func (s *AssessmentService) DeleteAkademikPenilaian(id string) error {
	return s.akademik.DeletePenilaian(id)
}

// ---- Jasmani ----

func (s *AssessmentService) GetJasmaniTree() ([]models.Jasmani, error) {
	return s.jasmani.GetAllComponents()
}

func (s *AssessmentService) GetJasmaniWithNilai(serdikID string) ([]models.Jasmani, map[uint64]models.PenilaianJasmani, error) {
	components, err := s.jasmani.GetAllComponents()
	if err != nil {
		return nil, nil, err
	}
	penilaian, err := s.jasmani.GetPenilaianBySerdik(serdikID)
	if err != nil {
		return nil, nil, err
	}
	nilaiMap := make(map[uint64]models.PenilaianJasmani)
	for _, p := range penilaian {
		nilaiMap[p.JasmaniComponentID] = p
	}
	return components, nilaiMap, nil
}

func (s *AssessmentService) CreateJasmaniPenilaian(p *models.PenilaianJasmani) error {
	return s.jasmani.CreatePenilaian(p)
}

func (s *AssessmentService) UpdateJasmaniPenilaian(id string, nilai float64, catatan *string) error {
	return s.jasmani.UpdatePenilaian(id, nilai, catatan)
}

func (s *AssessmentService) DeleteJasmaniPenilaian(id string) error {
	return s.jasmani.DeletePenilaian(id)
}

// ---- Mental ----

func (s *AssessmentService) GetMentalTree() ([]models.Mental, error) {
	return s.mental.GetAllComponents()
}

func (s *AssessmentService) GetMentalWithNilai(serdikID string) ([]models.Mental, map[uint64]models.PenilaianMental, error) {
	components, err := s.mental.GetAllComponents()
	if err != nil {
		return nil, nil, err
	}
	penilaian, err := s.mental.GetPenilaianBySerdik(serdikID)
	if err != nil {
		return nil, nil, err
	}
	nilaiMap := make(map[uint64]models.PenilaianMental)
	for _, p := range penilaian {
		nilaiMap[p.MentalComponentID] = p
	}
	return components, nilaiMap, nil
}

func (s *AssessmentService) CreateMentalPenilaian(p *models.PenilaianMental) error {
	return s.mental.CreatePenilaian(p)
}

func (s *AssessmentService) UpdateMentalPenilaian(id string, nilai float64, catatan *string) error {
	return s.mental.UpdatePenilaian(id, nilai, catatan)
}

func (s *AssessmentService) DeleteMentalPenilaian(id string) error {
	return s.mental.DeletePenilaian(id)
}

// ---- Sociometry ----

func (s *AssessmentService) GetPeriods(pokjarID *int) ([]models.SociometryPeriod, error) {
	return s.sociometry.GetPeriods(pokjarID)
}

func (s *AssessmentService) CreatePeriod(p *models.SociometryPeriod) error {
	return s.sociometry.CreatePeriod(p)
}

func (s *AssessmentService) ClosePeriod(id int64) error {
	return s.sociometry.ClosePeriod(id)
}

func (s *AssessmentService) SubmitEvaluation(e *models.SociometryEvaluation) error {
	return s.sociometry.CreateEvaluation(e)
}

// ---- Health ----

func (s *AssessmentService) GetHealthBySerdik(serdikID int64) (*models.SerdikHealthData, error) {
	return s.health.GetBySerdikID(serdikID)
}

func (s *AssessmentService) UpsertHealthData(d *models.SerdikHealthData) error {
	return s.health.UpsertHealthData(d)
}

func (s *AssessmentService) AddHealthRecord(r *models.HealthRecord) error {
	return s.health.CreateRecord(r)
}
