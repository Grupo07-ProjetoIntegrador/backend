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
		INSERT INTO lojas(luc, nome, segmento, status, email)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := database.DB.Exec(query, novaLoja.LUC, novaLoja.Nome, novaLoja.Segmento, novaLoja.Status, novaLoja.Email)

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
	var luc string

	nomeLimpo := strings.TrimSpace(nome)
	if nomeLimpo == "" {
		return "", fmt.Errorf("nome da loja nao informado")
	}

	// Verifica se a loja existe pelo nome, se o LUC existe e não está vazio, e puxa o ID
	queryBusca := `
		SELECT id, luc 
		FROM lojas 
		WHERE TRIM(UPPER(nome)) = TRIM(UPPER($1)) 
		  AND status = true 
		  AND luc IS NOT NULL 
		  AND TRIM(luc) <> ''
		LIMIT 1
	`

	err := database.DB.QueryRow(queryBusca, nomeLimpo).Scan(&lojaID, &luc)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("loja '%s' nao encontrada, inativa ou sem LUC cadastrado", nome)
		}

		return "", fmt.Errorf("erro inesperado ao buscar loja pelo nome: %v", err)
	}

	return lojaID, nil
}

// BuscarEmailLojaPorNome retorna o e-mail cadastrado para uma loja ativa com o nome informado.
func BuscarEmailLojaPorNome(nome string) (string, error) {
	var email string

	if strings.TrimSpace(nome) == "" {
		return "", fmt.Errorf("nome da loja nao informado")
	}

	queryBusca := `SELECT email FROM lojas WHERE nome = $1 AND status = true AND email IS NOT NULL AND TRIM(email) <> '' LIMIT 1`

	err := database.DB.QueryRow(queryBusca, nome).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("email nao encontrado para a loja")
		}

		return "", fmt.Errorf("erro inesperado ao buscar email da loja: %v", err)
	}

	return email, nil
}

// BuscarLojasComEmailPorSegmento lista as lojas ativas com e-mail cadastrado, opcionalmente filtrando por segmento.
func BuscarLojasComEmailPorSegmento(segmento string) ([]models.Loja, error) {
	query := `
		SELECT id, luc, nome, segmento, status, email
		FROM lojas
		WHERE status = true
		  AND email IS NOT NULL
		  AND TRIM(email) <> ''
	`

	args := []any{}
	seg := strings.TrimSpace(strings.ToLower(segmento))
	if seg != "" {
		if seg == "lojas" {
			// "Lojas" significa todas as lojas exceto Alimentação e Academia
			query += ` AND lower(segmento) NOT IN ('alimentação','academia','alimentacao')`
		} else {
			// filtro por igualdade (case-insensitive)
			query += ` AND lower(segmento) = $1`
			args = append(args, seg)
		}
	}

	query += ` ORDER BY nome ASC`

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar lojas com email: %v", err)
	}
	defer rows.Close()

	lojas := make([]models.Loja, 0)
	for rows.Next() {
		var loja models.Loja
		if err := rows.Scan(&loja.ID, &loja.LUC, &loja.Nome, &loja.Segmento, &loja.Status, &loja.Email); err != nil {
			return nil, fmt.Errorf("erro ao ler lojas com email: %v", err)
		}
		lojas = append(lojas, loja)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar lojas com email: %v", err)
	}

	return lojas, nil
}
