package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
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

func UpdateTreinamentosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	//id do treinamento existente
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "O id do treinamento é obrigatório", http.StatusBadRequest)
		return
	}

	//passa os dados do treinamento que quer modificar
	var UpdateTreinamento models.Treinamento

	//faz a leitura do JSON e retorna um erro caso haja
	err := json.NewDecoder(r.Body).Decode(&UpdateTreinamento)
	if err != nil {
		http.Error(w, "Erro ao fazer o update da ediçao", http.StatusBadRequest)
		return
	}

	//manda o id e o treinamento para a
	err = repositories.UpdateTreinamento(id, UpdateTreinamento)

	if err != nil {
		fmt.Println("Erro ao update do treinamento", err)
		http.Error(w, "Erro ao fazer Upload", http.StatusInternalServerError)
		return
	}

	//Resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensagem": "Treinamento editado com sucesso!"}`))
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

	treinamento, err := repositories.BuscarTreinamentoPorID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Treinamento não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar treinamento", http.StatusInternalServerError)
		return
	}

	payload := map[string]any{
		"treinamento_id": id,
		"treinamento": AutomacoesTreinamentoPayload{
			ID:            treinamento.ID,
			Tema:          treinamento.Tema,
			Descricao:     treinamento.Descricao,
			Objetivo:      treinamento.Objetivo,
			Data:          treinamento.Data,
			HorarioInicio: treinamento.HorarioInicio,
			HorarioFim:    treinamento.HorarioFim,
			Local:         treinamento.Local,
			SegmentoAlvo:  treinamento.SegmentoAlvo,
		},
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

	apiURL := automacoesBaseURL() + "/api/automacoes/gerar-forms"
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

// DispararConviteTreinamentoHandler envia convites segmentados para os destinatários selecionados.
func DispararConviteTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
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

	var req DisparoConviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao ler os dados do disparo", http.StatusBadRequest)
		return
	}

	if req.TreinamentoID == "" {
		http.Error(w, "O id de treinamento é obrigatório", http.StatusBadRequest)
		return
	}

	treinamento, err := repositories.BuscarTreinamentoPorID(req.TreinamentoID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Treinamento não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar treinamento", http.StatusInternalServerError)
		return
	}

	destinatarios, err := resolverDestinatariosDoDisparo(req, treinamento)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	modo := req.Modo
	if strings.TrimSpace(modo) == "" {
		modo = "individual"
	}

	payload := map[string]any{
		"treinamento_id": req.TreinamentoID,
		"modo":           modo,
		"treinamento": AutomacoesTreinamentoPayload{
			ID:            treinamento.ID,
			Tema:          treinamento.Tema,
			Descricao:     treinamento.Descricao,
			Objetivo:      treinamento.Objetivo,
			Data:          treinamento.Data,
			HorarioInicio: treinamento.HorarioInicio,
			HorarioFim:    treinamento.HorarioFim,
			Local:         treinamento.Local,
			SegmentoAlvo:  treinamento.SegmentoAlvo,
		},
		"destinatarios":        destinatarios,
		"segmento_loja":        req.SegmentoLoja,
		"segmento_treinamento": req.SegmentoTreinamento,
		"user_id":              req.UserID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Erro ao serializar payload", http.StatusInternalServerError)
		return
	}

	apiURL := automacoesBaseURL() + "/api/automacoes/disparar-convite"
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, "Erro ao chamar automacoes", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Erro ao disparar convites: %s", string(body)), http.StatusBadGateway)
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
	w.Write([]byte(`{"mensagem": "Treinamento editado com sucesso!"}`))
	w.Write([]byte(`{"mensagem":"Disparo de convites iniciado"}`))
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

	treinamento, err := repositories.BuscarTreinamentoPorID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Treinamento não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar treinamento", http.StatusInternalServerError)
		return
	}

	userID := r.URL.Query().Get("user_id")

	// Payload de exclusão
	deletePayload := map[string]any{
		"treinamento_id": id,
	}
	if userID != "" {
		deletePayload["user_id"] = userID
	}
	deleteJson, err := json.Marshal(deletePayload)
	if err != nil {
		http.Error(w, "Erro ao serializar payload de deleção", http.StatusInternalServerError)
		return
	}

	// Payload de geração
	generatePayload := map[string]any{
		"treinamento_id": id,
		"treinamento": AutomacoesTreinamentoPayload{
			ID:            treinamento.ID,
			Tema:          treinamento.Tema,
			Descricao:     treinamento.Descricao,
			Objetivo:      treinamento.Objetivo,
			Data:          treinamento.Data,
			HorarioInicio: treinamento.HorarioInicio,
			HorarioFim:    treinamento.HorarioFim,
			Local:         treinamento.Local,
			SegmentoAlvo:  treinamento.SegmentoAlvo,
		},
	}
	if userID != "" {
		generatePayload["user_id"] = userID
	}
	generateJson, err := json.Marshal(generatePayload)
	if err != nil {
		http.Error(w, "Erro ao serializar payload de geração", http.StatusInternalServerError)
		return
	}

	apiDeleteURL := "http://localhost:8000/api/automacoes/apagar-form"
	deleteResp, err := http.Post(apiDeleteURL, "application/json", bytes.NewBuffer(deleteJson))
	if err != nil {
		http.Error(w, "Erro ao chamar automacoes para apagar", http.StatusBadGateway)
		return
	}
	defer deleteResp.Body.Close()

	if deleteResp.StatusCode >= 400 && deleteResp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(deleteResp.Body)
		http.Error(w, fmt.Sprintf("Erro ao apagar formulario: %s", string(body)), http.StatusBadGateway)
		return
	}

	var deletePayloadResponse map[string]any
	if deleteResp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(deleteResp.Body)
		if len(body) > 0 {
			if err := json.Unmarshal(body, &deletePayloadResponse); err != nil {
				deletePayloadResponse = nil
			}
		}
	}

	apiGerarURL := "http://localhost:8000/api/automacoes/gerar-forms"
	resp, err := http.Post(apiGerarURL, "application/json", bytes.NewBuffer(generateJson))
	if err != nil {
		http.Error(w, "Erro ao chamar automacoes para gerar", http.StatusBadGateway)
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
	if deletePayloadResponse != nil {
		if value, ok := deletePayloadResponse["drive_deleted"]; ok {
			responsePayload["drive_deleted"] = value
		}
		if value, ok := deletePayloadResponse["form_id"]; ok {
			responsePayload["form_id"] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(responsePayload)
}
