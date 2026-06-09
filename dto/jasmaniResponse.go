package dto

type JasmaniResponse struct {
	ID       uint64             `json:"id"`
	Code     string             `json:"code"`
	Name     string             `json:"name"`
	Weight   float64            `json:"weight"`
	Nilai    *float64           `json:"nilai,omitempty"`
	Catatan  *string            `json:"catatan,omitempty"`
	Children []JasmaniResponse  `json:"children,omitempty"`
}