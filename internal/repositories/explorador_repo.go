// package repositories

// import (
// 	"fmt"

// 	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
// 	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
// )

// // BuscarExploradorLojas retorna todas as lojas ativas com métricas de treinamento
// // e taxa de presença calculadas dentro do período [dataInicio, dataFim].
// //
// // Modelagem das JOINs:
// //
// //	lojas (id, nome, luc, segmento, status)
// //	  └── presencas (loja_id, treinamento_id, status_presenca)
// //	          └── treinamentos (id, horario_inicio)
// //
// // ⚠ FIX #1: O filtro de data está no ON do LEFT JOIN com presencas,
// // NÃO no WHERE. Isso garante que lojas sem inscrições no período
// // ainda apareçam com totalTreinamentos = 0 e taxaParticipacao = 0.
// // Se o filtro ficasse no WHERE, o LEFT JOIN viraria um INNER JOIN efetivo.
// func BuscarExploradorLojas(dataInicio, dataFim string) ([]models.LojaExploradorItem, error) {
// 	query := `
// 		SELECT
// 			l.id,
// 			COALESCE(l.luc, '')         AS luc,
// 			COALESCE(l.nome, '')        AS nome,
// 			COALESCE(l.segmento, '')    AS segmento,

// 			-- Total de treinamentos distintos em que a loja foi inscrita no período
// 			COUNT(DISTINCT p.treinamento_id)                                          AS total_treinamentos,

// 			-- Presenças confirmadas (PRESENTE) dentro do período
// 			COUNT(CASE WHEN p.status_presenca = 'PRESENTE' THEN 1 END)               AS total_presentes,

// 			-- Total de inscrições válidas (PRESENTE + AUSENTE) dentro do período
// 			COUNT(CASE WHEN p.status_presenca IN ('PRESENTE', 'AUSENTE') THEN 1 END) AS total_inscricoes

// 		FROM lojas l

// 		-- ✅ O filtro de datas fica no ON — preserva o LEFT JOIN mesmo para lojas sem registros
// 		LEFT JOIN presencas p ON
// 			p.loja_id = l.id

// 		LEFT JOIN treinamentos t ON
// 			t.id = p.treinamento_id
// 			AND t.horario_inicio >= $1::timestamptz
// 			AND t.horario_inicio <= $2::timestamptz

// 		WHERE
// 			l.status = true
// 			AND l.email IS NOT NULL
// 			AND TRIM(l.email) <> ''

// 		GROUP BY l.id, l.luc, l.nome, l.segmento
// 		ORDER BY l.nome ASC
// 	`

// 	rows, err := database.DB.Query(query, dataInicio+"T00:00:00Z", dataFim+"T23:59:59Z")
// 	if err != nil {
// 		return nil, fmt.Errorf("erro ao buscar explorador de lojas: %v", err)
// 	}
// 	defer rows.Close()

// 	items := make([]models.LojaExploradorItem, 0)
// 	for rows.Next() {
// 		var item models.LojaExploradorItem
// 		var totalPresentes, totalInscricoes int

// 		// ✅ A ordem do Scan espelha exatamente a ordem do SELECT acima:
// 		// id, luc, nome, segmento, total_treinamentos, total_presentes, total_inscricoes
// 		if err := rows.Scan(
// 			&item.ID,
// 			&item.LUC,
// 			&item.Nome,
// 			&item.Segmento,
// 			&item.TotalTreinamentos,
// 			&totalPresentes,
// 			&totalInscricoes,
// 		); err != nil {
// 			return nil, fmt.Errorf("erro ao ler linha do explorador: %v", err)
// 		}

// 		// Calcula % de participação: Presentes / Total de Inscrições × 100
// 		if totalInscricoes > 0 {
// 			item.TaxaParticipacao = int((float64(totalPresentes) / float64(totalInscricoes)) * 100)
// 		} else {
// 			item.TaxaParticipacao = 0
// 		}

// 		items = append(items, item)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("erro ao iterar explorador de lojas: %v", err)
// 	}

// 	return items, nil
// }

package repositories

import (
	"fmt"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
)

// BuscarExploradorLojas retorna todas as lojas ativas com métricas de treinamento
// e taxa de presença calculadas dentro do período [dataInicio, dataFim].
//
// Modelagem das JOINs:
//
//  lojas (id, nome, luc, segmento, status)
//    └── treinamentos (id, horario_inicio) -> Filtrados pelo período
//         └── presencas (loja_id, treinamento_id, status_presenca)
//
// ⚠ FIX #1: O filtro de data foi reestruturado nos LEFT JOINs. Primeiro filtramos
// os treinamentos do período e depois vinculamos as presenças correspondentes.
// Isso garante que lojas sem inscrições no período ainda apareçam com totalTreinamentos = 0
// e taxaParticipacao = 0, sem inflar os dados com presenças de fora do período selecionado.
func BuscarExploradorLojas(dataInicio, dataFim string) ([]models.LojaExploradorItem, error) {
	query := `
		SELECT
			l.id,
			COALESCE(l.luc, '')         AS luc,
			COALESCE(l.nome, '')        AS nome,
			COALESCE(l.segmento, '')    AS segmento,

			-- Total de treinamentos distintos em que a loja foi inscrita no período filtrado
			COUNT(DISTINCT p.treinamento_id)                                            AS total_treinamentos,

			-- Presenças confirmadas (PRESENTE) dentro do período filtrado
			COUNT(CASE WHEN p.status_presenca = 'PRESENTE' THEN 1 END)                 AS total_presentes,

			-- Total de inscrições válidas (PRESENTE + AUSENTE) dentro do período filtrado
			COUNT(CASE WHEN p.status_presenca IN ('PRESENTE', 'AUSENTE') THEN 1 END)   AS total_inscricoes

		FROM lojas l

		-- 🟢 Primeiro filtramos os treinamentos que aconteceram no período desejado
		LEFT JOIN treinamentos t ON
			t.horario_inicio >= $1::timestamptz
			AND t.horario_inicio <= $2::timestamptz

		-- 🟢 Só puxamos as presenças que pertencem aos treinamentos desse período específico
		LEFT JOIN presencas p ON
			p.loja_id = l.id
			AND p.treinamento_id = t.id

		WHERE
			l.status = true
			AND l.email IS NOT NULL
			AND TRIM(l.email) <> ''

		GROUP BY l.id, l.luc, l.nome, l.segmento
		ORDER BY l.nome ASC
	`

	rows, err := database.DB.Query(query, dataInicio+"T00:00:00Z", dataFim+"T23:59:59Z")
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar explorador de lojas: %v", err)
	}
	defer rows.Close()

	items := make([]models.LojaExploradorItem, 0)
	for rows.Next() {
		var item models.LojaExploradorItem
		var totalPresentes, totalInscricoes int

		// ✅ A ordem do Scan espelha exatamente a ordem do SELECT acima:
		// id, luc, nome, segmento, total_treinamentos, total_presentes, total_inscricoes
		if err := rows.Scan(
			&item.ID,
			&item.LUC,
			&item.Nome,
			&item.Segmento,
			&item.TotalTreinamentos,
			&totalPresentes,
			&totalInscricoes,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler linha do explorador: %v", err)
		}

		// Calcula % de participação: Presentes / Total de Inscrições × 100
		if totalInscricoes > 0 {
			item.TaxaParticipacao = int((float64(totalPresentes) / float64(totalInscricoes)) * 100)
		} else {
			item.TaxaParticipacao = 0
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar explorador de lojas: %v", err)
	}

	return items, nil
}