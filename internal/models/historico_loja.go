package models

// TreinamentoLoja representa um treinamento agregado pelo ponto de vista de uma loja específica.
// Retornado pelo endpoint GET /api/lojas/:id/historico
type TreinamentoLojaItem struct {
	TreinamentoID string   `json:"treinamento_id"`
	Tema          string   `json:"tema"`
	Data          string   `json:"data"`           // YYYY-MM-DD
	HorarioInicio string   `json:"horario_inicio"` // HH:MM
	Presentes     []string `json:"presentes"`
	Ausentes      []string `json:"ausentes"`
}
