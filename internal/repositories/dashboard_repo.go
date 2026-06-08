package repositories

import (
	"database/sql"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
)

func ObterDadosDashboard() (models.DashboardStats, error) {
	var stats models.DashboardStats

	// 1. Total de Participações reais (Apenas status válidos)
	queryParticipacoes := `SELECT COUNT(*) FROM presencas WHERE status_presenca IN ('PRESENTE', 'PENDENTE')`
	err := database.DB.QueryRow(queryParticipacoes).Scan(&stats.TotalParticipacoes)
	if err != nil {
		return stats, err
	}

	// 2. Total de Lojas Impactadas REAIS (Ignora lojas genéricas automáticas)
	queryLojasImpactadas := `
		SELECT COUNT(DISTINCT p.loja_id) 
		FROM presencas p
		INNER JOIN lojas l ON p.loja_id = l.id
		WHERE p.status_presenca = 'PRESENTE' 
		  AND l.segmento <> 'Não Informado'
	`
	err = database.DB.QueryRow(queryLojasImpactadas).Scan(&stats.TotalLojasImpactadas)
	if err != nil {
		return stats, err
	}

	// 3. Total de Lojas Cadastradas Oficiais (Totalmente idêntica à do Supabase)
	queryTotalLojas := `SELECT COUNT(*) FROM lojas WHERE email IS NOT NULL AND TRIM(email) <> '' AND status = true`
	err = database.DB.QueryRow(queryTotalLojas).Scan(&stats.TotalLojasCadastradas)
	if err != nil {
		return stats, err
	}

	// // 3. Total de Lojas Cadastradas Oficiais (Desconsidera inserções automáticas de teste)
	// queryTotalLojas := `SELECT COUNT(*) FROM lojas WHERE segmento <> 'Não Informado'`
	// err = database.DB.QueryRow(queryTotalLojas).Scan(&stats.TotalLojasCadastradas)
	// if err != nil {
	// 	return stats, err
	// }

	// 4. Taxa de Presença Média
	var capacidadeTotal int
	queryCapacidade := `
		SELECT COALESCE(SUM(capacidade_maxima), 0) 
		FROM treinamentos 
		WHERE id IN (SELECT DISTINCT treinamento_id FROM presencas)
	`
	err = database.DB.QueryRow(queryCapacidade).Scan(&capacidadeTotal)
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

	queryTreinamentosComPresenca := `
		SELECT COUNT(*)
		FROM (
			SELECT DISTINCT p.loja_id, p.treinamento_id
			FROM presencas p
			INNER JOIN lojas l ON p.loja_id = l.id
			WHERE p.status_presenca = 'PRESENTE'
			  AND l.status = true
			  AND l.segmento <> 'Não Informado'
		) presencas_por_treinamento
	`
	err = database.DB.QueryRow(queryTreinamentosComPresenca).Scan(&stats.TotalTreinamentosComPresenca)
	if err != nil {
		return stats, err
	}

	// 6. Top Engajamento (Lojas oficiais mais presentes)
	queryTop := `
		SELECT l.nome, COUNT(DISTINCT p.treinamento_id) as total
		FROM presencas p
		INNER JOIN lojas l ON p.loja_id = l.id
		WHERE p.status_presenca = 'PRESENTE'
		  AND l.status = true
		  AND l.segmento <> 'Não Informado'
		GROUP BY l.id, l.nome
		ORDER BY total DESC, l.nome ASC
		LIMIT 5
	`
	rowsTop, err := database.DB.Query(queryTop)
	if err == nil {
		defer rowsTop.Close()
		for rowsTop.Next() {
			var r models.LojaRanking
			if err := rowsTop.Scan(&r.Name, &r.Total); err == nil {
				stats.TopEngajamento = append(stats.TopEngajamento, r)
			}
		}
	}

	// 7. Radar de Risco (Garante valores padrão mesmo se os registros forem nulos)
	queryRisco := `
		SELECT 
			l.nome, 
			l.luc,
			COALESCE(COUNT(CASE WHEN p.status_presenca = 'AUSENTE' THEN 1 END), 0) as total_faltas,
			COALESCE(TO_CHAR(MAX(CASE WHEN p.status_presenca = 'PRESENTE' THEN p.data_registro END), 'DD/MM/YYYY'), 'Nunca') as ultima_data
		FROM lojas l
		LEFT JOIN presencas p ON l.id = p.loja_id
		WHERE l.status = true AND l.email IS NOT NULL AND TRIM(l.email) <> ''
		GROUP BY l.nome, l.luc
		ORDER BY total_faltas DESC, l.nome ASC
		LIMIT 5
	`
	rowsRisco, err := database.DB.Query(queryRisco)
	if err == nil {
		defer rowsRisco.Close()
		for rowsRisco.Next() {
			var r models.LojaRanking
			var lucString sql.NullString

			// Scan seguro tratando strings e inteiros nulos do banco de dados
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

	// 8. Gráfico de Evolução Mensal
	meses := []string{"Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"}
	for _, m := range meses {
		stats.EvolucaoMensal = append(stats.EvolucaoMensal, models.MensalStat{Mes: m, Taxa: stats.TaxaPresencaMedia})
	}

	return stats, nil
}
