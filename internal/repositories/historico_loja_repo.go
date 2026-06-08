package repositories

import (
	"fmt"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
)

// BuscarHistoricoLoja retorna todos os treinamentos em que uma loja participou,
// agrupando os participantes por treinamento em arrays "presentes" e "ausentes".
//
// A query faz JOIN entre presencas e treinamentos filtrando por loja_id.
// O filtro de data é opcional — se dataInicio e dataFim forem strings vazias, retorna tudo.
func BuscarHistoricoLoja(lojaID, dataInicio, dataFim string) ([]models.TreinamentoLojaItem, error) {

	// Monta a cláusula de filtro de data dinamicamente
	filtroData := ""
	args := []interface{}{lojaID}

	if dataInicio != "" && dataFim != "" {
		filtroData = `AND t.horario_inicio >= $2::timestamptz AND t.horario_inicio <= $3::timestamptz`
		args = append(args, dataInicio+"T00:00:00Z", dataFim+"T23:59:59Z")
	}

	query := fmt.Sprintf(`
		SELECT
			p.treinamento_id,
			COALESCE(t.tema, 'Treinamento sem Tema')                        AS tema,
			TO_CHAR(t.horario_inicio, 'YYYY-MM-DD')                         AS data,
			TO_CHAR(t.horario_inicio, 'HH24:MI')                            AS horario_inicio,
			p.nome_participante,
			p.status_presenca
		FROM presencas p
		INNER JOIN treinamentos t ON t.id = p.treinamento_id
		WHERE p.loja_id = $1
		%s
		ORDER BY t.horario_inicio DESC, p.nome_participante ASC
	`, filtroData)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar histórico da loja: %v", err)
	}
	defer rows.Close()

	// Agrupa por treinamento_id para montar os arrays presentes/ausentes
	order := []string{} // mantém a ordem de inserção
	itemMap := map[string]*models.TreinamentoLojaItem{}

	for rows.Next() {
		var (
			treinamentoID string
			tema          string
			data          string
			horario       string
			nomePartic    string
			status        string
		)

		if err := rows.Scan(&treinamentoID, &tema, &data, &horario, &nomePartic, &status); err != nil {
			continue
		}

		if _, exists := itemMap[treinamentoID]; !exists {
			itemMap[treinamentoID] = &models.TreinamentoLojaItem{
				TreinamentoID: treinamentoID,
				Tema:          tema,
				Data:          data,
				HorarioInicio: horario,
				Presentes:     []string{},
				Ausentes:      []string{},
			}
			order = append(order, treinamentoID)
		}

		item := itemMap[treinamentoID]
		switch status {
		case "PRESENTE", "CONFIRMADO", "SIM":
			item.Presentes = append(item.Presentes, nomePartic)
		default:
			item.Ausentes = append(item.Ausentes, nomePartic)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar histórico da loja: %v", err)
	}

	// Retorna na ordem original (mais recente primeiro, conforme ORDER BY)
	result := make([]models.TreinamentoLojaItem, 0, len(order))
	for _, id := range order {
		result = append(result, *itemMap[id])
	}

	return result, nil
}
