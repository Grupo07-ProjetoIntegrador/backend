package repositories

import (
	"database/sql"
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
		SELECT id, tema, segmento_alvo, horario_inicio, conteudo, status, capacidade_maxima
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

		err := linhas.Scan(&t.ID, &t.Tema, &t.Segmento, &dataHoraBanco, &t.Conteudo, &t.Status, &t.CapacidadeMaxima)

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

// DeletarTreinamento remove um treinamento do banco de dados usando o seu ID

func DeletarTreinamento(id string) error {
	//Escrever o comando de delete do SQL

	query := `DELETE FROM treinamentos WHERE id = $1`
	fmt.Println("DEBUG: O ID que chegou no banco é:", id)

	//Comando Exec vai ate o banco de dados e pega o valor do id

	resultado, err := database.DB.Exec(query, id)

	//Verifica se a conexão com banco de dados falhou
	if err != nil {
		return err
	}

	//Pergunta quantas linhas foram apagadas com esse comando

	linhasAfetadas, err := resultado.RowsAffected()
	fmt.Println("DEBUG: Quantidade de linhas que o Supabase apagou:", linhasAfetadas)
	if err != nil {
		return err
	}

	//Caso nenhuma linha seja apagada significa que nao existe treinamento com aquele id fornecido

	if linhasAfetadas == 0 {
		return fmt.Errorf("Nenhum treinamento encontrado com esse ID")
	}

	//Caso nao tenha nenhum erro vai retornar o nil

	return nil

}

// BuscarTreinamentoTema retorna o tema do treinamento pelo ID
func BuscarTreinamentoTema(id string) (string, error) {
	var tema string

	query := `SELECT tema FROM treinamentos WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(&tema)
	if err != nil {
		return "", err
	}

	return tema, nil
}

// BuscarTreinamentoPorID retorna os dados completos do treinamento para automacoes.
func BuscarTreinamentoPorID(id string) (models.Treinamento, error) {
	var t models.Treinamento

	query := `
		SELECT
			COALESCE(id::text, ''),
			COALESCE(tema, ''),
			COALESCE(descricao, ''),
			COALESCE(categoria, ''),
			COALESCE(TO_CHAR(data, 'YYYY-MM-DD'), ''),
			COALESCE(TO_CHAR(horario_inicio, 'HH24:MI'), ''),
			COALESCE(TO_CHAR(horario_fim, 'HH24:MI'), ''),
			COALESCE(local, ''),
			COALESCE(modalidade, ''),
			COALESCE(conteudo, ''),
			COALESCE(capacidade_maxima, 0),
			COALESCE(segmento_alvo, ''),
			COALESCE(status::text, ''),
			COALESCE(objetivo, ''),
			COALESCE(observacoes, ''),
			COALESCE(material_apoio, ''),
			COALESCE(responsavel, ''),
			COALESCE(area_responsavel, ''),
			COALESCE(tags, ''),
			COALESCE(recorrente, false)
		FROM treinamentos
		WHERE id = $1
		LIMIT 1
	`

	err := database.DB.QueryRow(query, id).Scan(
		&t.ID,
		&t.Tema,
		&t.Descricao,
		&t.Categoria,
		&t.Data,
		&t.HorarioInicio,
		&t.HorarioFim,
		&t.Local,
		&t.Modalidade,
		&t.Conteudo,
		&t.CapacidadeMaxima,
		&t.SegmentoAlvo,
		&t.Status,
		&t.Objetivo,
		&t.Observacoes,
		&t.MaterialApoio,
		&t.Responsavel,
		&t.AreaResponsavel,
		&t.Tags,
		&t.Recorrente,
	)
	if err != nil {
		return models.Treinamento{}, err
	}

	return t, nil
}

// BuscarFormularioTreinamento retorna o link do formulario e o id do Google Form
func BuscarFormularioTreinamento(id string) (string, string, error) {
	var url string
	var formID string

	query := `
		SELECT url_formulario, google_form_id
		FROM formularios_treinamento
		WHERE treinamento_id = $1
		ORDER BY criado_em DESC
		LIMIT 1
	`

	err := database.DB.QueryRow(query, id).Scan(&url, &formID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", sql.ErrNoRows
		}
		return "", "", err
	}

	return url, formID, nil
}

// ListarTreinamentos retorna todos os treinamentos do banco ordenados por data
// func ListarTreinamentos() ([]models.Treinamento, error) {
// 	query := `
// 		SELECT
// 			COALESCE(id::text, ''),
// 			COALESCE(tema, ''),
// 			COALESCE(descricao, ''),
// 			COALESCE(categoria, ''),
// 			COALESCE(data::text, ''),
// 			COALESCE(horario_inicio::text, ''),
// 			COALESCE(horario_fim::text, ''),
// 			COALESCE(local, ''),
// 			COALESCE(modalidade, ''),
// 			COALESCE(conteudo, ''),
// 			COALESCE(capacidade_maxima, 0),
// 			COALESCE(segmento_alvo, ''),
// 			COALESCE(status::text, ''),
// 			COALESCE(objetivo, ''),
// 			COALESCE(observacoes, ''),
// 			COALESCE(material_apoio, ''),
// 			COALESCE(responsavel, ''),
// 			COALESCE(area_responsavel, ''),
// 			COALESCE(tags, ''),
// 			COALESCE(recorrente, false)
// 		FROM treinamentos
// 		ORDER BY data DESC
// 	`

// 	rows, err := database.DB.Query(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("erro ao listar treinamentos: %v", err)
// 	}
// 	defer rows.Close()

// 	var treinamentos []models.Treinamento
// 	for rows.Next() {
// 		var t models.Treinamento
// 		err := rows.Scan(
// 			&t.ID, &t.Tema, &t.Descricao, &t.Categoria, &t.Data,
// 			&t.HorarioInicio, &t.HorarioFim, &t.Local, &t.Modalidade, &t.Conteudo,
// 			&t.CapacidadeMaxima, &t.SegmentoAlvo, &t.Status,
// 			&t.Objetivo, &t.Observacoes, &t.MaterialApoio,
// 			&t.Responsavel, &t.AreaResponsavel, &t.Tags, &t.Recorrente,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("erro ao ler treinamento do banco: %v", err)
// 		}
// 		treinamentos = append(treinamentos, t)
// 	}

// 	if treinamentos == nil {
// 		treinamentos = []models.Treinamento{}
// 	}

// 	return treinamentos, nil
// >>>>>>> automacoes-python
// }
