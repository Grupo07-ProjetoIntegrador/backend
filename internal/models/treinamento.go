package models

// Treinamento representa o molde exato da tabela 'treinamentos' no banco de dados
type Treinamento struct {
	ID               string `json:"id"`                // No banco é 'uuid', no Go usamos 'string'
	Tema             string `json:"tema"`              // 'varchar'
	Descricao        string `json:"descricao"`         // 'text'
	Categoria        string `json:"categoria"`         // 'varchar'
	Data             string `json:"data"`              // 'timestamp' no banco, 'string' facilita receber o JSON do front
	HorarioInicio    string `json:"horario_inicio"`    // 'timestamp'
	HorarioFim       string `json:"horario_fim"`       // 'timestamp'
	Local            string `json:"local"`             // 'text'
	Modalidade       string `json:"modalidade"`        // 'varchar'
	Conteudo         string `json:"conteudo"`          // 'text'
	CapacidadeMaxima int    `json:"capacidade_maxima"` // 'int4' no banco, traduzido para 'int' no Go
	SegmentoAlvo     string `json:"segmento_alvo"`     // 'varchar'
	Status           string `json:"status"`            // 'Treinamento_status' (Enum no banco, 'string' no Go)

	Objetivo      string `json:"objetivo"`
	Observacoes   string `json:"observacoes"`
	MaterialApoio string `json:"material_apoio"`

	Responsavel     string `json:"responsavel"`
	AreaResponsavel string `json:"area_responsavel"`
	Tags            string `json:"tags"`
	Recorrente      bool   `json:"recorrente"`
	LocalID         string `json:"local_id"`
}

type TreinamentoResumo struct {
	ID               string `json:"id"`
	Tema             string `json:"tema"`
	Segmento         string `json:"segmento"`
	Data             string `json:"data"`
	DataHora         string `json:"data_hora"`
	HorarioInicio    string `json:"horario_inicio"`
	Conteudo         string `json:"conteudo"`
	Status           string `json:"status"`
	CapacidadeMaxima int    `json:"capacidade_maxima"`
}
