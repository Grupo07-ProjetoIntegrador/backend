package models

// Presenca representa o molde exato da tabela 'presencas' no banco de dados
type Presenca struct {
	ID               string `json:"id"`                // uuid no banco
	TreinamentoID    string `json:"treinamento_id"`    // uuid no banco (Chave estrangeira)
	LojaID           string `json:"loja_id"`           // uuid no banco (Chave estrangeira)
	NomeParticipante string `json:"nome_participante"` // varchar
	StatusPresenca   string `json:"status_presenca"`   // presenca_status (Enum)
	Email            string `json:"email"`             // varchar
	Telefone         string `json:"telefone"`          // varchar
	Cargo            string `json:"cargo"`             // varchar
	DataRegistro     string `json:"data_registro"`     // timestamp
}

type PresencaResponse struct {
	ID            string `json:"id"`
	LUC           string `json:"luc"`
	Loja          string `json:"loja"`
	Representante string `json:"representante"`
	Status        string `json:"status"`
}

type CriarPresencaInput struct {
	TreinamentoID string `json:"treinamento_id"`
	LUC           string `json:"luc"`
	Representante string `json:"representante"`
	Status        string `json:"status"`
}
