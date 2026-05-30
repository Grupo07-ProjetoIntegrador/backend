package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	// Importando as pastas do seu projeto
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

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

// CadastrarTreinamentoHandler recebe os dados da tela "Cadastrar Novo Treinamento"
func CadastrarTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
	// Liberar o CORS para o Front-end conseguir acessar (Mantendo a versão completa do seu colega)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Se o navegador estiver apenas testando a conexão (Preflight OPTIONS), retorna OK
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
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Treinamento '%s' criado com sucesso! O ID para o Google Forms é: %s", novoTreinamento.Tema, idGerado)

	// Geração do formulário agora é manual.
}

// ListarTreinamentosHandler busca os dados do treinamento e lança na tela de lista
func ListarTreinamentosHandler(w http.ResponseWriter, r *http.Request) {
	// Liberar o CORS para o Front-end conseguir acessar
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Se o navegador estiver apenas testando a conexão (Preflight OPTIONS), retorna OK
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Verificando se esta usando o comando Get
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	// Buscando a lista do banco de dados
	lista, err := repositories.ListarTreinamentos()

	// Verificacao de erro de conexao
	if err != nil {
		http.Error(w, "Erro ao buscar a lista de treinamentos", http.StatusInternalServerError)
		fmt.Println("Erro na listagem:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)
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

// GerarFormularioTreinamentoHandler dispara a geracao manual do Google Forms
func GerarFormularioTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "O id de treinamento é obrigatório", http.StatusBadRequest)
		return
	}

	tema, err := repositories.BuscarTreinamentoTema(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Treinamento não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar treinamento", http.StatusInternalServerError)
		return
	}

	payload := map[string]string{
		"treinamento_id": id,
		"tema":           tema,
	}
	// Propagar user_id se fornecido (para que o serviço de automacoes use as credenciais do usuario)
	userID := r.URL.Query().Get("user_id")
	if userID != "" {
		payload["user_id"] = userID
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Erro ao serializar payload", http.StatusInternalServerError)
		return
	}

	apiURL := "http://localhost:8000/api/automacoes/gerar-forms"
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, "Erro ao chamar automacoes", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Erro ao gerar formulario: %s", string(body)), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"mensagem":"Geração do formulário iniciada"}`))
}

// BuscarFormularioTreinamentoHandler retorna o link do formulario quando existir
func BuscarFormularioTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
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

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "O id de treinamento é obrigatório", http.StatusBadRequest)
		return
	}

	urlFormulario, formID, err := repositories.BuscarFormularioTreinamento(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Formulario ainda nao gerado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar formulario", http.StatusInternalServerError)
		return
	}

	if formID == "" && urlFormulario != "" {
		re := regexp.MustCompile(`/forms/d/(?:e/)?([^/]+)`)
		matches := re.FindStringSubmatch(urlFormulario)
		if len(matches) > 1 {
			formID = matches[1]
		}
	}

	editURL := ""
	if formID != "" {
		editURL = fmt.Sprintf("https://docs.google.com/forms/d/%s/edit", formID)
	}

	creatorName, creatorEmail := resolverCriadorFormulario(urlFormulario)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url_formulario": urlFormulario,
		"url_edicao":     editURL,
		"creator_name":   creatorName,
		"creator_email":  creatorEmail,
	})
}

// ApagarFormularioTreinamentoHandler remove o vinculo do formulario no backend
func ApagarFormularioTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido. Use DELETE.", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "O id de treinamento é obrigatório", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("user_id")

	payload := map[string]string{
		"treinamento_id": id,
	}
	if userID != "" {
		payload["user_id"] = userID
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Erro ao serializar payload", http.StatusInternalServerError)
		return
	}

	apiURL := "http://localhost:8000/api/automacoes/apagar-form"
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, "Erro ao chamar automacoes", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "Formulario ainda nao gerado", http.StatusNotFound)
		return
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Erro ao apagar formulario: %s", string(body)), http.StatusBadGateway)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	if len(body) > 0 {
		w.WriteHeader(http.StatusOK)
		w.Write(body)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensagem":"Formulario removido"}`))
}

// RegerarFormularioTreinamentoHandler apaga o formulario e cria um novo
func RegerarFormularioTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "O id de treinamento é obrigatório", http.StatusBadRequest)
		return
	}

	tema, err := repositories.BuscarTreinamentoTema(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Treinamento não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar treinamento", http.StatusInternalServerError)
		return
	}

	payload := map[string]string{
		"treinamento_id": id,
		"tema":           tema,
	}
	userID := r.URL.Query().Get("user_id")
	if userID != "" {
		payload["user_id"] = userID
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Erro ao serializar payload", http.StatusInternalServerError)
		return
	}

	apiDeleteURL := "http://localhost:8000/api/automacoes/apagar-form"
	deleteResp, err := http.Post(apiDeleteURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, "Erro ao chamar automacoes", http.StatusBadGateway)
		return
	}
	defer deleteResp.Body.Close()

	if deleteResp.StatusCode >= 400 && deleteResp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(deleteResp.Body)
		http.Error(w, fmt.Sprintf("Erro ao apagar formulario: %s", string(body)), http.StatusBadGateway)
		return
	}

	var deletePayload map[string]any
	if deleteResp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(deleteResp.Body)
		if len(body) > 0 {
			if err := json.Unmarshal(body, &deletePayload); err != nil {
				deletePayload = nil
			}
		}
	}

	apiGerarURL := "http://localhost:8000/api/automacoes/gerar-forms"
	resp, err := http.Post(apiGerarURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, "Erro ao chamar automacoes", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Erro ao gerar formulario: %s", string(body)), http.StatusBadGateway)
		return
	}

	responsePayload := map[string]any{
		"mensagem": "Geração do formulário iniciada",
	}
	if deletePayload != nil {
		if value, ok := deletePayload["drive_deleted"]; ok {
			responsePayload["drive_deleted"] = value
		}
		if value, ok := deletePayload["form_id"]; ok {
			responsePayload["form_id"] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(responsePayload)
}
