package repositories

import (
	"fmt"
	"strings"
	"time"

	// Importa a conexão com o banco
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"

	// Importa a Struct 'Treinamento' que você acabou de criar
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
)

// InserirTreinamento salva um novo evento no banco e devolve o ID (UUID) gerado
func InserirTreinamento(t models.Treinamento) (string, error) {
	// Variável que vai guardar o UUID gerado pelo banco
	var idGerado string

	//Tratamento(conversão do formato da data para string)
	dataParseada, err := time.Parse("02/01/2006", t.Data)

	if err != nil {
		return "", fmt.Errorf("erro ao converter a data para o formato dd/mm/aaaa: %v", err)
	}

	//Formata para o padrao do banco de dados
	dataFormatadaBanco := dataParseada.Format("2006-01-02")

	//Tratamento da HorarioInicio para ser aceito no banco.
	inicioParseado, err := time.Parse("02/01/2006 15:04", t.Data+" "+t.HorarioInicio)
	if err != nil {
		return "", fmt.Errorf("Erro ao converter horario inicio: %v", err)
	}

	horarioInicioBanco := inicioParseado.Format("2006-01-02 15:04:00")

	//Tratamento do HorarioFim para ser aceito no banco de dados

	FimParseado, err := time.Parse("02/01/2006 15:04", t.Data+" "+t.HorarioFim)

	if err != nil {
		return "", fmt.Errorf("Erro ao converter horario fim: %v", err)
	}

	horarioFimBanco := FimParseado.Format("2006-01-02 15:04:00")

	//Tratamento do Status para fazer a string ficar toda Maiuscula para colocar no banco de dados
	statusBanco := strings.ToUpper(t.Status)

	// O 'RETURNING id' no final é o segredo para o banco devolver o UUID criado na mesma hora.
	query := `
		INSERT INTO treinamentos (
			tema, descricao, categoria, data, horario_inicio, 
			horario_fim, local, modalidade, conteudo, 
			capacidade_maxima, segmento_alvo, status,
			objetivo, observacoes, material_apoio,
			responsavel, area_responsavel, tags, recorrente
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id
	`

	// O QueryRow executa o comando e o Scan(&idGerado) pega a resposta do banco e salva na nossa variável
	err = database.DB.QueryRow(
		query,
		t.Tema, t.Descricao, t.Categoria, dataFormatadaBanco, horarioInicioBanco,
		horarioFimBanco, t.Local, t.Modalidade, t.Conteudo,
		t.CapacidadeMaxima, t.SegmentoAlvo, statusBanco,
		t.Objetivo, t.Observacoes, t.MaterialApoio,
		t.Responsavel, t.AreaResponsavel, t.Tags, t.Recorrente,
	).Scan(&idGerado)

	if err != nil {
		// Se der erro, retorna uma string vazia e o erro para o Handler
		return "", fmt.Errorf("erro ao inserir o treinamento: %v", err)
	}

	// Se deu tudo certo, retorna o ID novinho em folha e "nil" para o erro!
	return idGerado, nil
}

func ListarTreinamentos() ([]models.TreinamentoResumo, error) {
	//Começando Montando a query para selecionar as informações do banco de dados
	query := `
		SELECT id, tema, segmento_alvo, horario_inicio, conteudo, status
		FROM treinamentos
		ORDER BY horario_inicio DESC
	`
	//Utilizando a query e acessando o banco de dados
	linhas, err := database.DB.Query(query)

	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar treinamentos: %v", err)

	}

	defer linhas.Close()
	// Criando lista para guardar os dados que veem do banco de dados
	var lista []models.TreinamentoResumo

	//Loop para percorrer cada linha que o banco de dados devolveu
	for linhas.Next() {
		var t models.TreinamentoResumo

		var dataHoraBanco time.Time

		err := linhas.Scan(&t.ID, &t.Tema, &t.Segmento, &dataHoraBanco, &t.Conteudo, &t.Status)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler os dados da linha: %v", err)
		}

		t.Data = dataHoraBanco.Format("02 Jan 2006 às 15:04")

		t.DataHora = dataHoraBanco.Format("2006-01-02T15:04:00")
		t.HorarioInicio = dataHoraBanco.Format("15:04")

		lista = append(lista, t)
	}

	return lista, nil

}
