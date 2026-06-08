package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

// HistoricoLojaHandler responde a GET /api/lojas/historico?id=<uuid>&data_inicio=YYYY-MM-DD&data_fim=YYYY-MM-DD
//
// Retorna todos os treinamentos em que a loja participou no período,
// agrupados com arrays "presentes" e "ausentes" por treinamento.
func HistoricoLojaHandler(w http.ResponseWriter, r *http.Request) {
	// Só aceita GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Cabeçalhos CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Parâmetros obrigatórios / opcionais
	lojaID := strings.TrimSpace(r.URL.Query().Get("id"))
	if lojaID == "" {
		http.Error(w, `{"error":"parâmetro 'id' é obrigatório"}`, http.StatusBadRequest)
		return
	}

	dataInicio := strings.TrimSpace(r.URL.Query().Get("data_inicio"))
	dataFim := strings.TrimSpace(r.URL.Query().Get("data_fim"))

	historico, err := repositories.BuscarHistoricoLoja(lojaID, dataInicio, dataFim)
	if err != nil {
		http.Error(w, `{"error":"Erro ao buscar histórico da loja"}`, http.StatusInternalServerError)
		return
	}

	// Garante JSON array vazio em vez de null
	if historico == nil {
		historico = []models.TreinamentoLojaItem{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(historico)
}
