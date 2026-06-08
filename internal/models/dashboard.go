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
	Mes  string `json:"name"`
	Taxa int    `json:"Taxa de Presença"`
}

type LojaRanking struct {
	Name           string `json:"name"`
	Total          int    `json:"total"`
	LUC            string `json:"luc,omitempty"`
	Faltas         int    `json:"faltas"`
	UltimaPresenca string `json:"ultimaPresenca"`
}
