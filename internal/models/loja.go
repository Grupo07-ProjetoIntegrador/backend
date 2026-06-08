package models

type Loja struct {
	ID       string `json:"id"`
	LUC      string `json:"luc"`
	Nome     string `json:"nome"`
	Segmento string `json:"segmento"`
	Status   bool   `json:"status"`
	Email    string `json:"email,omitempty"`
}

// LojaExploradorItem representa uma loja com métricas calculadas para o Explorador de Lojas.
// Retornado pelo endpoint GET /api/lojas/explorador?data_inicio=YYYY-MM-DD&data_fim=YYYY-MM-DD
type LojaExploradorItem struct {
	ID                string `json:"id"`
	LUC               string `json:"luc"`
	Nome              string `json:"nome"`
	Segmento          string `json:"segmento"`
	TotalTreinamentos int    `json:"totalTreinamentos"`
	TaxaParticipacao  int    `json:"taxaParticipacao"` // 0-100 (percentual)
}
