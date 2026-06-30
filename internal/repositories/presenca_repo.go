package repositories

import (
	"database/sql"
	"fmt"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
)

// InserirPresencaPendente salva a inscrição vinda do Forms com o status inicial 'PENDENTE'
func InserirPresencaPendente(treinamentoID string, lojaID string, nomeParticipante string, email string, telefone string, cargo string) error {
	// 1. O COMANDO SQL
	// Baseado na estrutura da sua tabela, inserimos os dados e fixamos o status
	query := `
        INSERT INTO presencas (treinamento_id, loja_id, nome_participante, email, telefone, cargo, status_presenca) 
        VALUES ($1, $2, $3, $4, $5, $6, 'PENDENTE')
    `

	// 2. A EXECUÇÃO NO BANCO
	// Substituímos os parâmetros pelas variáveis que o handler nos enviou
	_, err := database.DB.Exec(query, treinamentoID, lojaID, nomeParticipante, email, telefone, cargo)

	// 3. TRATAMENTO DE ERRO
	if err != nil {
		// Se houver algum problema (ex: o ID do treinamento não existir no banco),
		// o Supabase vai bloquear e nós devolvemos esse erro para o Handler.
		return fmt.Errorf("erro ao inserir presença pendente no banco: %v", err)
	}

	return nil
}

// AtualizarStatusPresenca recebe o ID do registro de presença e o novo status (ex: "PRESENTE" ou "AUSENTE")
// Essa função será usada no dia do evento quando o admin clicar no sistema para confirmar a participação.
func AtualizarStatusPresenca(presencaID string, novoStatus string) error {
	// 1. Comando SQL para ATUALIZAR apenas o status de uma presença específica
	query := `
        UPDATE presencas 
        SET status_presenca = $1 
        WHERE id = $2
    `

	// 2. Executa o comando no Supabase substituindo $1 pelo novo status e $2 pelo ID
	_, err := database.DB.Exec(query, novoStatus, presencaID)

	// 3. Tratamento de erro
	if err != nil {
		return fmt.Errorf("erro ao atualizar o status da presença: %v", err)
	}

	return nil
}

// SalvarPresencaPlanilha insere os dados vindos de importação de planilhas (Mantido do seu colega)
func SalvarPresencaPlanilha(treinamentoID string, luc string, nomeParticipante string, status string) error {
	// 2. A Mágica da Subquery no SQL
	// Na hora de inserir o loja_id, nós fazemos um (SELECT id FROM lojas WHERE luc = $2).
	// Omitimos email, telefone e cargo porque a planilha antiga não tem esses dados.

	query := `
        INSERT INTO presencas (treinamento_id, loja_id, nome_participante, status_presenca)
        VALUES(
            $1,
            (SELECT id FROM lojas WHERE luc = $2 LIMIT 1),
            $3,
            $4
        )
    `
	// Executa a query
	_, err := database.DB.Exec(query, treinamentoID, luc, nomeParticipante, status)

	if err != nil {
		return fmt.Errorf("erro ao inserir presença do LUC %s: %v", luc, err)
	}

	return nil
} // <- Chave de fechamento que estava faltando!

func ListarPresencaPorTreinamentos(treinamentoID string) ([]models.PresencaResponse, error) {

	presencas := []models.PresencaResponse{}

	query := `
            SELECT p.id, l.luc, l.nome, p.nome_participante, p.status_presenca
            FROM presencas p
            INNER JOIN lojas l ON p.loja_id = l.id
            WHERE p.treinamento_id = $1
            ORDER BY p.data_registro DESC
    `

	rows, err := database.DB.Query(query, treinamentoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.PresencaResponse
		err := rows.Scan(&p.ID, &p.LUC, &p.Loja, &p.Representante, &p.Status)
		if err != nil {
			continue
		}
		presencas = append(presencas, p)
	}

	// Removemos aquela checagem "if presencas == nil" do final porque ela não é mais necessária!

	return presencas, nil
}

// ConfirmarPresencaPorEmail atualiza o status_presenca de 'PENDENTE' para 'PRESENTE'
// buscando pela combinação de treinamento_id e email que esteja com status 'PENDENTE'. (Mantido do seu código)
func ConfirmarPresencaPorEmail(treinamentoID string, email string) error {
	query := `
        UPDATE presencas 
        SET status_presenca = 'PRESENTE' 
        WHERE treinamento_id = $1 AND email = $2 AND status_presenca = 'PENDENTE'
    `

	res, err := database.DB.Exec(query, treinamentoID, email)
	if err != nil {
		return fmt.Errorf("erro ao atualizar o status de presença para PRESENTE: %v", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("nenhuma inscrição PENDENTE encontrada para este e-mail neste treinamento")
	}

	return nil
}

// BuscarNomeParticipantePorEmail retorna o nome cadastrado na presença pendente.
func BuscarNomeParticipantePorEmail(treinamentoID string, email string) (string, error) {
	var nome string

	query := `
		SELECT nome_participante
		FROM presencas
		WHERE treinamento_id = $1 AND email = $2
		ORDER BY id DESC
		LIMIT 1
	`

	err := database.DB.QueryRow(query, treinamentoID, email).Scan(&nome)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("participante nao encontrado para este e-mail neste treinamento")
		}

		return "", fmt.Errorf("erro ao buscar nome do participante: %v", err)
	}

	return nome, nil
}

// CriarPresencaManual insere uma nova presença direto no banco de dados
func CriarPresencaManual(treinamentoID string, luc string, nomeParticipante string, status string) error {
	var lojaID string

	// 1. Busca o ID da loja baseado no LUC enviado pelo Front
	queryLoja := `SELECT id FROM lojas WHERE luc = $1 LIMIT 1`
	err := database.DB.QueryRow(queryLoja, luc).Scan(&lojaID)
	if err != nil {
		return fmt.Errorf("loja com o LUC %s não foi encontrada no sistema", luc)
	}

	// 2. Insere o participante na tabela de presenças
	queryPresenca := `
		INSERT INTO presencas (treinamento_id, loja_id, nome_participante, status_presenca, data_registro)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err = database.DB.Exec(queryPresenca, treinamentoID, lojaID, nomeParticipante, status)
	if err != nil {
		return fmt.Errorf("erro ao inserir presença no banco: %v", err)
	}

	return nil
}

// DeletarPresenca remove um registro de presença do banco de dados pelo ID
func DeletarPresenca(id string) error {
	query := `DELETE FROM presencas WHERE id = $1`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar presença do banco: %v", err)
	}

	// Verifica se alguma linha realmente foi afetada (se o ID existia)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("nenhum registro de presença encontrado com o ID fornecido")
	}

	return nil
}

// EditarPresenca atualiza o LUC (loja_id), nome do participante e status no banco de dados pelo ID da presenca
func EditarPresenca(id string, luc string, nomeParticipante string, status string) error {
	var lojaID string

	// 1. Busca o ID da loja baseado no LUC enviado pelo Front
	queryLoja := `SELECT id FROM lojas WHERE luc = $1 LIMIT 1`
	err := database.DB.QueryRow(queryLoja, luc).Scan(&lojaID)
	if err != nil {
		return fmt.Errorf("loja com o LUC %s não foi encontrada no sistema", luc)
	}

	// 2. Atualiza os dados na tabela de presencas
	query := `
		UPDATE presencas 
		SET loja_id = $1, nome_participante = $2, status_presenca = $3
		WHERE id = $4
	`
	_, err = database.DB.Exec(query, lojaID, nomeParticipante, status, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar presença no banco: %v", err)
	}

	return nil
}
