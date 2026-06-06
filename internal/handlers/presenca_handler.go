package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

type ConfirmarPresencaRequest struct {
	TreinamentoID string `json:"treinamento_id"`
	Email         string `json:"email"`
}

func ListarPresencasHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Configuração do CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	treinamentoID := r.URL.Query().Get("treinamento_id")
	if treinamentoID == "" {
		http.Error(w, `{"erro": "ID do treinamento é obrigatório"}`, http.StatusBadRequest)
		return
	}

	presencas, err := repositories.ListarPresencaPorTreinamentos(treinamentoID)

	if err != nil {
		// Se o erro for apenas porque não há linhas, tratamos como sucesso com lista vazia
		// Você pode checar se o erro é 'sql.ErrNoRows' ou se prefere apenas zerar o erro se o seu grupo preferir
		fmt.Printf("🚨 Erro ou aviso ao buscar no banco: %v\n", err)

		// Se quiser que mesmo com erro ele não quebre o front, podemos forçar o envio de uma lista vazia:
		presencas = []models.PresencaResponse{}
	}

	// Se a busca deu certo mas veio nula (sem registros no banco), transformamos em [] para o React não quebrar
	if presencas == nil {
		presencas = []models.PresencaResponse{} // Evita mandar 'null' no JSON, manda '[]'
	}

	// 3. Devolve os dados em formato JSON para o React (Sempre 200 OK se chegou aqui)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(presencas)
}

// ConfirmarPresencaHandler atende o PATCH enviado pelo checkin.html
func ConfirmarPresencaHandler(w http.ResponseWriter, r *http.Request) {
	// CORS Headers - Permite que a página HTML executada no celular faça chamadas para esta API
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Responde requisições OPTIONS pré-voo do CORS
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, "Método não permitido. Use PATCH.", http.StatusMethodNotAllowed)
		return
	}

	var req ConfirmarPresencaRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Erro ao ler os dados enviados", http.StatusBadRequest)
		return
	}

	if req.TreinamentoID == "" || req.Email == "" {
		http.Error(w, "Campos 'treinamento_id' e 'email' são obrigatórios", http.StatusBadRequest)
		return
	}

	nomeParticipante, _ := repositories.BuscarNomeParticipantePorEmail(req.TreinamentoID, req.Email)

	// Tenta atualizar a presença para PRESENTE
	err = repositories.ConfirmarPresencaPorEmail(req.TreinamentoID, req.Email)
	if err != nil {
		// Retorna 404 se não achar a inscrição pendente ou outro erro correspondente
		http.Error(w, fmt.Sprintf("Erro ao confirmar presença: %v", err), http.StatusNotFound)
		return
	}

	if err := notificarPresencaValidada(req.TreinamentoID, req.Email, nomeParticipante); err != nil {
		fmt.Println("Aviso ao enviar e-mail de presença validada:", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Presença confirmada com sucesso para o participante %s!", req.Email)
}

func notificarPresencaValidada(treinamentoID string, email string, nomeParticipante string) error {
	treinamento, err := repositories.BuscarTreinamentoPorID(treinamentoID)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"treinamento_id": treinamentoID,
		"treinamento": map[string]any{
			"id":             treinamento.ID,
			"tema":           treinamento.Tema,
			"descricao":      treinamento.Descricao,
			"objetivo":       treinamento.Objetivo,
			"data":           treinamento.Data,
			"horario_inicio": treinamento.HorarioInicio,
			"horario_fim":    treinamento.HorarioFim,
			"local":          treinamento.Local,
			"segmento_alvo":  treinamento.SegmentoAlvo,
		},
		"destinatario": map[string]any{
			"nome":  nomeParticipante,
			"email": email,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	apiURL := os.Getenv("AUTOMACOES_PUBLIC_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8000"
	}

	resp, err := http.Post(apiURL+"/api/automacoes/notificar-presenca-validada", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("automacao retornou status %d: %s", resp.StatusCode, string(responseBody))
	}

	return nil
}

func CriarPresencaManualHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"erro": "Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}

	var input models.CriarPresencaInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, `{"erro": "JSON inválido"}`, http.StatusBadRequest)
		return
	}

	// Validação dos campos obrigatórios
	if input.TreinamentoID == "" || input.LUC == "" || input.Representante == "" {
		http.Error(w, `{"erro": "Todos os campos são obrigatórios"}`, http.StatusBadRequest)
		return
	}

	input.Status = strings.ToUpper(input.Status)

	// Chama a função do repositório enviando os dados validados
	err = repositories.CriarPresencaManual(input.TreinamentoID, input.LUC, input.Representante, input.Status)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Participante adicionado com sucesso!"})
}

func DeletarPresencaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	// Aceita apenas o método DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, `{"erro": "Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}

	// Pega o ID dos parâmetros da URL (?id=...)
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"erro": "O parâmetro ID é obrigatório"}`, http.StatusMethodNotAllowed)
		return
	}

	// Chama a função do repositório para deletar do banco
	err := repositories.DeletarPresenca(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Presença removida com sucesso!"})
}
