package repositories

import (
	"fmt"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
)

// InserirPresencaPendente salva a inscrição vinda do Forms com o status inicial 'PENDENTE'
func InserirPresencaPendente(treinamentoID string, lojaID string, nomeParticipante string) error {
	// 1. O COMANDO SQL
	// Baseado na estrutura da sua tabela, inserimos os dados e fixamos o status
	query := `
		INSERT INTO presencas (treinamento_id, loja_id, nome_participante, status_presenca) 
		VALUES ($1, $2, $3, 'PENDENTE')
	`

	// 2. A EXECUÇÃO NO BANCO
	// Substituímos o $1, $2 e $3 pelas variáveis que o handler nos enviou
	_, err := database.DB.Exec(query, treinamentoID, lojaID, nomeParticipante)

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

//Função para salvar os dados de planilhas

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
	//Executa a query
	_, err := database.DB.Exec(query, treinamentoID, luc, nomeParticipante, status)

	if err != nil {
		return fmt.Errorf("erro ao inserir presença do LUC %s: %v", luc, err)
	}

	return nil

}
