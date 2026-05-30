package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
)

func InserirLoja(novaLoja models.Loja) error {
	//Comandos SQL
	query := `
		INSERT INTO lojas(luc, nome, segmento, status)
		VALUES ($1, $2, $3, $4)
	`
	_, err := database.DB.Exec(query, novaLoja.LUC, novaLoja.Nome, novaLoja.Segmento, novaLoja.Status)

	if err != nil {
		return fmt.Errorf("Erro ao inserir loja no banco: %v", err)
	}

	return nil
}

// BuscarOuCriarLoja procura uma loja pelo LUC. Se não achar, cria e devolve o novo ID.
func BuscarOuCriarLoja(luc string, nome string) (string, error) {
	var lojaID string

	// 1. A BUSCA (SELECT)
	// Primeiro, perguntamos ao banco: "Você tem alguma loja com este LUC?"
	queryBusca := `SELECT id FROM lojas WHERE luc = $1`

	err := database.DB.QueryRow(queryBusca, luc).Scan(&lojaID)

	if err != nil {

		if err == sql.ErrNoRows {

			queryCriacao := `
				INSERT INTO lojas (luc, nome, segmento, status) 
				VALUES ($1, $2, 'Não Informado', true) 
				RETURNING id
			`

			errInsert := database.DB.QueryRow(queryCriacao, luc, nome).Scan(&lojaID)

			if errInsert != nil {
				return "", fmt.Errorf("falha ao criar nova loja pelo webhook: %v", errInsert)
			}

			return lojaID, nil
		}

		// Se for qualquer outro erro (ex: banco fora do ar), interrompe e avisa
		return "", fmt.Errorf("erro inesperado ao buscar loja pelo LUC: %v", err)
	}

	return lojaID, nil
}

// BuscarLojaAtivaPorNome procura uma loja ativa pelo nome. Nao cria loja automaticamente.
func BuscarLojaAtivaPorNome(nome string) (string, error) {
	var lojaID string

	if strings.TrimSpace(nome) == "" {
		return "", fmt.Errorf("nome da loja nao informado")
	}

	queryBusca := `SELECT id FROM lojas WHERE nome = $1 AND status = true`

	err := database.DB.QueryRow(queryBusca, nome).Scan(&lojaID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("loja nao encontrada ou inativa")
		}

		return "", fmt.Errorf("erro inesperado ao buscar loja pelo nome: %v", err)
	}

	return lojaID, nil
}
