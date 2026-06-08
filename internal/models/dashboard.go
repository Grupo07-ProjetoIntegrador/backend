package models

type DashboardStats struct {
	TotalParticipacoes           int           `json:"totalParticipacoes"`
	TotalTreinamentosComPresenca int           `json:"totalTreinamentosComPresenca"`
	TotalLojasImpactadas         int           `json:"totalLojasImpactadas"`
	TotalLojasCadastradas        int           `json:"totalLojasCadastradas"`
	TaxaPresencaMedia            int           `json:"taxaPresencaMedia"`
	MediaColabPorLoja            float64       `json:"mediaColabPorLoja"`
	EvolucaoMensal               []MensalStat  `json:"evolucaoMensal"`
	TopEngajamento               []LojaRanking `json:"topEngajamento"`
	RadarRisco                   []LojaRanking `json:"radarRisco"`
}

type MensalStat struct {
	Mes       string `json:"name"`      // Mudado para "name" para casar direto com o dataKey do Recharts
	Inscritos int    `json:"inscritos"` // Barra cinza do gráfico
	Taxa      int    `json:"taxa"`      // Se ainda usar a taxa antiga
	Presentes int    `json:"presentes"` // Barra vermelha do gráfico
}

type LojaRanking struct {
	Name           string `json:"name"`
	Total          int    `json:"total"`
	LUC            string `json:"luc,omitempty"`
	Faltas         int    `json:"faltas"`
	UltimaPresenca string `json:"ultimaPresenca"`
}

// DadosMensais representa a pontuação de um mês no gráfico
type DadosMensais struct {
	Name      string `json:"name"`      // "Jan", "Fev", "Mar"...
	Inscritos int    `json:"inscritos"` // Total de alunos na lista (capacidade máxima ou total de registros)
	Presentes int    `json:"presentes"` // Total de alunos com status "Presente"
}

// DashboardResponse estende a resposta que você já deve ter para os cards
type DashboardResponse struct {
	TaxaPresencaMedia     float64        `json:"taxaPresencaMedia"`
	TotalLojasImpactadas  int            `json:"totalLojasImpactadas"`
	TotalLojasCadastradas int            `json:"totalLojasCadastradas"`
	MediaColabPorLoja     float64        `json:"mediaColabPorLoja"`
	TotalParticipacoes    int            `json:"totalParticipacoes"`
	TopEngajamento        []interface{}  `json:"topEngajamento"`
	RadarRisco            []interface{}  `json:"radarRisco"`
	EvolucaoParticipacao  []DadosMensais `json:"evolucaoParticipacao"` // 🚀 NOSSO GRÁFICO AQUI
}
