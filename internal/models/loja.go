package models

type Loja struct {
	ID       string `json:"id"`
	LUC      string `json:"luc"`
	Nome     string `json:"nome"`
	Segmento string `json:"segmento"`
	Status   bool   `json:"status"`
	Email    string `json:"email,omitempty"`
}
