package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

func DashboardStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `{"erro": "Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}

	dados, err := repositories.ObterDadosDashboard()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"erro": "%v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dados)
}
