package repositories

import (
	"fmt"

	// Importa a conexão com o banco
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"

	// Importa a Struct 'Treinamento' que você acabou de criar
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
)

// InserirTreinamento salva um novo evento no banco e devolve o ID (UUID) gerado
func InserirTreinamento(t models.Treinamento) (string, error) {
	// Variável que vai guardar o UUID gerado pelo banco
	var idGerado string

	// O 'RETURNING id' no final é o segredo para o banco devolver o UUID criado na mesma hora.
	query := `
		INSERT INTO treinamentos (
			tema, descricao, categoria, data, horario_inicio, 
			horario_fim, local, modalidade, conteudo, 
			capacidade_maxima, segmento_alvo, status
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	// O QueryRow executa o comando e o Scan(&idGerado) pega a resposta do banco e salva na nossa variável
	err := database.DB.QueryRow(
		query,
		t.Tema, t.Descricao, t.Categoria, t.Data, t.HorarioInicio,
		t.HorarioFim, t.Local, t.Modalidade, t.Conteudo,
		t.CapacidadeMaxima, t.SegmentoAlvo, t.Status,
	).Scan(&idGerado)

	if err != nil {
		// Se der erro, retorna uma string vazia e o erro para o Handler
		return "", fmt.Errorf("erro ao inserir o treinamento: %v", err)
	}

	// Se deu tudo certo, retorna o ID novinho em folha e "nil" para o erro!
	return idGerado, nil
}
