package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	// Tenta atualizar a presença para PRESENTE
	err = repositories.ConfirmarPresencaPorEmail(req.TreinamentoID, req.Email)
	if err != nil {
		// Retorna 404 se não achar a inscrição pendente ou outro erro correspondente
		http.Error(w, fmt.Sprintf("Erro ao confirmar presença: %v", err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Presença confirmada com sucesso para o participante %s!", req.Email)
}
