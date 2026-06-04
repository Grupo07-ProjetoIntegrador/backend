package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	// Importando as pastas do seu projeto
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

// ListarTreinamentosHandler retorna todos os treinamentos em JSON
func ListarTreinamentosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	treinamentos, err := repositories.ListarTreinamentos()
	if err != nil {
		fmt.Println("Erro ao listar treinamentos:", err)
		http.Error(w, "Erro ao buscar treinamentos no banco de dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(treinamentos)
}

type AutomacoesTreinamentoPayload struct {
	ID            string `json:"id"`
	Tema          string `json:"tema"`
	Descricao     string `json:"descricao"`
	Objetivo      string `json:"objetivo"`
	Data          string `json:"data"`
	HorarioInicio string `json:"horario_inicio"`
	HorarioFim    string `json:"horario_fim"`
	Local         string `json:"local"`
	SegmentoAlvo  string `json:"segmento_alvo"`
}

type ConviteDestinatario struct {
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Segmento string `json:"segmento"`
}

type DisparoConviteRequest struct {
	TreinamentoID       string                `json:"treinamento_id"`
	Modo                string                `json:"modo"`
	SegmentoLoja        string                `json:"segmento_loja"`
	SegmentoTreinamento string                `json:"segmento_treinamento"`
	Destinatarios       []ConviteDestinatario `json:"destinatarios"`
	UserID              string                `json:"user_id"`
}

func resolverCriadorFormulario(urlFormulario string) (string, string) {
	if urlFormulario == "" {
		return "", ""
	}

	parsed, err := url.Parse(urlFormulario)
	if err != nil || parsed.Fragment == "" {
		return "", ""
	}

	fragmentValues, err := url.ParseQuery(parsed.Fragment)
	if err != nil {
		return "", ""
	}

	ownerID := fragmentValues.Get("owner_user_id")
	if ownerID == "" {
		return "", ""
	}

	var displayName string
	var email string
	err = database.DB.QueryRow(
		`SELECT display_name, email FROM profiles WHERE user_id = $1`,
		ownerID,
	).Scan(&displayName, &email)
	if err != nil {
		return "", ""
	}

	return displayName, email
}

func automacoesBaseURL() string {
	if baseURL := os.Getenv("AUTOMACOES_PUBLIC_URL"); baseURL != "" {
		return baseURL
	}

	return "http://localhost:8000"
}

func resolverDestinatariosDoDisparo(req DisparoConviteRequest, treinamento models.Treinamento) ([]ConviteDestinatario, error) {
	if len(req.Destinatarios) > 0 {
		resolved := make([]ConviteDestinatario, 0, len(req.Destinatarios))
		for _, destinatario := range req.Destinatarios {
			item := destinatario
			if strings.TrimSpace(item.Email) == "" && strings.TrimSpace(item.Nome) != "" {
				if email, err := repositories.BuscarEmailLojaPorNome(item.Nome); err == nil {
					item.Email = email
				}
			}

			if strings.TrimSpace(item.Email) == "" {
				return nil, fmt.Errorf("destinatário sem e-mail: %s", item.Nome)
			}

			resolved = append(resolved, item)
		}

		return resolved, nil
	}

	segmentoFiltro := ""
	switch req.Modo {
	case "segmento_treinamento":
		segmentoFiltro = treinamento.SegmentoAlvo
	case "segmento_loja":
		segmentoFiltro = req.SegmentoLoja
	}

	lojas, err := repositories.BuscarLojasComEmailPorSegmento(segmentoFiltro)
	if err != nil {
		return nil, err
	}
	if len(lojas) == 0 {
		return nil, fmt.Errorf("nenhuma loja com e-mail encontrado para o disparo")
	}

	destinatarios := make([]ConviteDestinatario, 0, len(lojas))
	for _, loja := range lojas {
		destinatarios = append(destinatarios, ConviteDestinatario{
			Nome:     loja.Nome,
			Email:    loja.Email,
			Segmento: loja.Segmento,
		})
	}

	return destinatarios, nil
}

// CadastrarTreinamentoHandler recebe os dados da tela "Cadastrar Novo Treinamento"
func CadastrarTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 1. Verifica se o Front-end está mandando um POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	// 2. Cria o "molde" vazio usando a Struct que você fez na Fase 1
	var novoTreinamento models.Treinamento

	// 3. Lê o JSON que veio do Front-end e preenche o molde
	err := json.NewDecoder(r.Body).Decode(&novoTreinamento)
	if err != nil {
		http.Error(w, "Erro ao ler os dados do formulário preenchido", http.StatusBadRequest)
		return
	}

	// 4. INTEGRAÇÃO COM O BANCO DE DADOS
	// Chama a função do repositório que salva e devolve o ID gerado
	idGerado, err := repositories.InserirTreinamento(novoTreinamento)

	if err != nil {
		// Se der erro (ex: banco offline), avisa o front-end
		http.Error(w, "Erro ao salvar o treinamento no banco de dados", http.StatusInternalServerError)
		fmt.Println("Erro no repositório:", err)
		return
	}

	// 5. RESPOSTA DE SUCESSO!
	// Retorna Status 201 (Criado) e devolve o UUID para o administrador copiar e usar no Google Forms
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Treinamento '%s' criado com sucesso! O ID para o Google Forms é: %s", novoTreinamento.Tema, idGerado)

	// Geração do formulário agora é manual.
}

// DeletarTreinamentoHandler remove um treinamento do banco através do ID
func DeletarTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
	// Configuracao do CORS para o front acessar
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Trava que faz a URL aceitar somente o DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido. Use DELETE.", http.StatusMethodNotAllowed)
		return
	}

	// Extrai o ID da URL
	id := r.URL.Query().Get("id")

	// Verifica se o ID foi enviado
	if id == "" {
		http.Error(w, "O id de treinamento é obrigatório", http.StatusBadRequest)
		return
	}

	// Chama a função do repositório
	err := repositories.DeletarTreinamento(id)

	// Verifica se tem algum erro
	if err != nil {
		http.Error(w, "Erro ao deletar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Caso dê tudo certo vem pra cá
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensagem": "Treinamento deletado com sucesso!"}`))
}
