package dto

type AkademikNodeResponse struct {
	ID                  uint64                  `json:"id"`
	AkademikComponentID uint64                  `json:"akademik_component_id"`
	Code                string                  `json:"code"`
	Name                string                  `json:"name"`
	Weight              float64                 `json:"weight"`
	Nilai               *float64                `json:"nilai,omitempty"`
	Catatan             *string                 `json:"catatan,omitempty"`
	Children            []AkademikNodeResponse  `json:"children,omitempty"`
}