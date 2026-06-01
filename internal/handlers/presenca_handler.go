package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

type ConfirmarPresencaRequest struct {
	TreinamentoID string `json:"treinamento_id"`
	Email         string `json:"email"`
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
