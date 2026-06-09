package dto

type ComponentTreeResponse struct {
	ID       uint64                   `json:"id"`
	ParentID *uint64                  `json:"parent_id"`
	Code     string                   `json:"code"`
	Name     string                   `json:"name"`
	Weight   float64                  `json:"weight"`
	Nilai    float64                  `json:"nilai"`
	Catatan  *string                  `json:"catatan,omitempty"`
	Children []ComponentTreeResponse  `json:"children"`
}