package models

// Presenca representa o molde exato da tabela 'presencas' no banco de dados
type Presenca struct {
	ID               string `json:"id"`                // uuid no banco
	TreinamentoID    string `json:"treinamento_id"`    // uuid no banco (Chave estrangeira)
	LojaID           string `json:"loja_id"`           // uuid no banco (Chave estrangeira)
	NomeParticipante string `json:"nome_participante"` // varchar
	StatusPresenca   string `json:"status_presenca"`   // presenca_status (Enum)
	DataRegistro     string `json:"data_registro"`     // timestamp
}
