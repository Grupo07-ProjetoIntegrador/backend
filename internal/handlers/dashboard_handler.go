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

	// 1. Captura os Query Params enviados pelo React na URL (?data_inicio=...&data_fim=...)
	dataInicio := r.URL.Query().Get("data_inicio")
	dataFim := r.URL.Query().Get("data_fim")

	// Fallbacks de segurança caso venham vazios (ex: primeiro carregamento da tela)
	if dataInicio == "" { dataInicio = "2026-01-01" }
	if dataFim == "" { dataFim = "2026-12-31" }

	// 2. CORREÇÃO: Passa as duas datas extraídas como argumentos
	dados, err := repositories.ObterDadosDashboard(dataInicio, dataFim)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"erro": "%v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dados)
}