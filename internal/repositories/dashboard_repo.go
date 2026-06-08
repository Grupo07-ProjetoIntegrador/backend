package repositories

import (
	"database/sql"
	"fmt"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
)

func ObterDadosDashboard(dataInicio, dataFim string) (models.DashboardStats, error) {
	var stats models.DashboardStats

	fmt.Println("======= FILTRO RECEBIDO NO BACKEND =======")
	fmt.Printf("Data Início: %s | Data Fim: %s\n", dataInicio, dataFim)
	fmt.Println("==========================================")

	if dataInicio == "" { dataInicio = "2026-01-01" }
	if dataFim == "" { dataFim = "2026-12-31" }

	// 1. Total de Participações reais
	queryParticipacoes := `
		SELECT COUNT(p.id) 
		FROM presencas p
		INNER JOIN treinamentos t ON p.treinamento_id = t.id
		WHERE p.status_presenca IN ('PRESENTE', 'PENDENTE')
		  AND t.data::date BETWEEN $1::date AND $2::date 
		  AND t.status != 'CANCELADO'
	`
	err := database.DB.QueryRow(queryParticipacoes, dataInicio, dataFim).Scan(&stats.TotalParticipacoes)
	if err != nil {
		return stats, err
	}

	// 2. Total de Lojas Impactadas REAIS
	queryLojasImpactadas := `
		SELECT COUNT(DISTINCT p.loja_id) 
		FROM presencas p
		INNER JOIN lojas l ON p.loja_id = l.id
		INNER JOIN treinamentos t ON p.treinamento_id = t.id
		WHERE p.status_presenca = 'PRESENTE' 
		  AND l.segmento <> 'Não Informado'
		  AND t.data::date BETWEEN $1::date AND $2::date 
		  AND t.status != 'CANCELADO'
	`
	err = database.DB.QueryRow(queryLojasImpactadas, dataInicio, dataFim).Scan(&stats.TotalLojasImpactadas)
	if err != nil {
		return stats, err
	}

	// 3. Total de Lojas Cadastradas Oficiais (Total geral mantido estático)
	queryTotalLojas := `SELECT COUNT(*) FROM lojas WHERE email IS NOT NULL AND TRIM(email) <> '' AND status = true`
	err = database.DB.QueryRow(queryTotalLojas).Scan(&stats.TotalLojasCadastradas)
	if err != nil {
		return stats, err
	}

	// 4. Taxa de Presença Média
	var capacidadeTotal int
	queryCapacidade := `
		SELECT COALESCE(SUM(capacidade_maxima), 0) 
		FROM treinamentos t
		WHERE t.data::date BETWEEN $1::date AND $2::date 
		  AND t.status != 'CANCELADO'
		  AND t.id IN (SELECT DISTINCT treinamento_id FROM presencas)
	`
	err = database.DB.QueryRow(queryCapacidade, dataInicio, dataFim).Scan(&capacidadeTotal)
	if err != nil {
		return stats, err
	}

	if capacidadeTotal > 0 {
		stats.TaxaPresencaMedia = (stats.TotalParticipacoes * 100) / capacidadeTotal
	}

	// 5. Média de Colaboradores por Loja Ativa
	if stats.TotalLojasImpactadas > 0 {
		stats.MediaColabPorLoja = float64(stats.TotalParticipacoes) / float64(stats.TotalLojasImpactadas)
	}

	// 6. Top Engajamento
	queryTop := `
		SELECT l.nome, COUNT(p.id) as total
		FROM presencas p
		INNER JOIN lojas l ON p.loja_id = l.id
		INNER JOIN treinamentos t ON p.treinamento_id = t.id
		WHERE p.status_presenca = 'PRESENTE' AND l.segmento <> 'Não Informado'
		  AND t.data::date BETWEEN $1::date AND $2::date 
		  AND t.status != 'CANCELADO'
		GROUP BY l.nome
		ORDER BY total DESC
		LIMIT 5
	`
	rowsTop, err := database.DB.Query(queryTop, dataInicio, dataFim)
	if err == nil {
		defer rowsTop.Close()
		for rowsTop.Next() {
			var r models.LojaRanking
			if err := rowsTop.Scan(&r.Name, &r.Total); err == nil {
				stats.TopEngajamento = append(stats.TopEngajamento, r)
			}
		}
	}

	// 7. Radar de Risco
	queryRisco := `
		SELECT 
			l.nome, 
			l.luc,
			COALESCE(COUNT(CASE WHEN p.status_presenca = 'AUSENTE' AND t.data::date BETWEEN $1::date AND $2::date THEN 1 END), 0) as total_faltas,
			COALESCE(TO_CHAR(MAX(CASE WHEN p.status_presenca = 'PRESENTE' AND t.data::date BETWEEN $1::date AND $2::date THEN p.data_registro END), 'DD/MM/YYYY'), 'Nunca') as ultima_data
		FROM lojas l
		LEFT JOIN presencas p ON l.id = p.loja_id
		LEFT JOIN treinamentos t ON p.treinamento_id = t.id
		WHERE l.status = true AND l.email IS NOT NULL AND TRIM(l.email) <> ''
		GROUP BY l.nome, l.luc
		ORDER BY total_faltas DESC, l.nome ASC
		LIMIT 5
	`
	rowsRisco, err := database.DB.Query(queryRisco, dataInicio, dataFim)
	if err == nil {
		defer rowsRisco.Close()
		for rowsRisco.Next() {
			var r models.LojaRanking
			var lucString sql.NullString

			errScan := rowsRisco.Scan(&r.Name, &lucString, &r.Faltas, &r.UltimaPresenca)
			if errScan == nil {
				if lucString.Valid {
					r.LUC = lucString.String
				} else {
					r.LUC = "Não informado"
				}
				stats.RadarRisco = append(stats.RadarRisco, r)
			}
		}
	}

	// 8. Gráfico de Evolução Mensal (Cronológico Dinâmico para viradas de ano)
	queryEvolucao := `
		SELECT 
			TO_CHAR(t.data, 'YYYY-MM') AS ano_mes,
			TO_CHAR(t.data, 'Mon') AS mes_nome,
			COUNT(p.id) AS total_inscritos,
			COUNT(CASE WHEN UPPER(TRIM(p.status_presenca::text)) = 'PRESENTE' THEN 1 END) AS total_presentes
		FROM treinamentos t
		LEFT JOIN presencas p ON t.id = p.treinamento_id
		WHERE t.data::date BETWEEN $1::date AND $2::date 
		  AND UPPER(TRIM(t.status::text)) != 'CANCELADO'
		GROUP BY ano_mes, mes_nome
		ORDER BY ano_mes ASC;
	`

	rowsEvolucao, err := database.DB.Query(queryEvolucao, dataInicio, dataFim)
	
	// Reinicializa o array dinâmico
	stats.EvolucaoMensal = []models.MensalStat{}

	if err == nil {
		defer rowsEvolucao.Close()

		traducoes := map[string]string{
			"Jan": "Jan", "Feb": "Fev", "Mar": "Mar", "Apr": "Abr",
			"May": "Mai", "Jun": "Jun", "Jul": "Jul", "Aug": "Ago",
			"Sep": "Set", "Oct": "Out", "Nov": "Nov", "Dec": "Dez",
		}

		for rowsEvolucao.Next() {
			var anoMes, mesNome string
			var inscritos, presentes int

			errScan := rowsEvolucao.Scan(&anoMes, &mesNome, &inscritos, &presentes)
			if errScan == nil {
				nomeTraduzido := mesNome
				if brNome, ok := traducoes[mesNome]; ok {
					nomeTraduzido = brNome
				}

				taxa := 0
				if inscritos > 0 {
					taxa = (presentes * 100) / inscritos
				}

				// Alimenta o array mantendo rigorosamente a ordem sequencial temporal retornada pelo banco
				stats.EvolucaoMensal = append(stats.EvolucaoMensal, models.MensalStat{
					Mes:       nomeTraduzido,
					Inscritos: inscritos,
					Presentes: presentes,
					Taxa:      taxa,
				})
			}
		}
	}

	// Fallback estrutural: Se o período customizado selecionado não tiver nenhuma linha de dados cadastrada
	if len(stats.EvolucaoMensal) == 0 {
		mesesPadrao := []string{"Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"}
		for _, m := range mesesPadrao {
			stats.EvolucaoMensal = append(stats.EvolucaoMensal, models.MensalStat{Mes: m, Taxa: 0, Inscritos: 0, Presentes: 0})
		}
	}

	return stats, nil
}