package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

// ExploradorLojasHandler atende GET /api/lojas/explorador
//
// Query Params:
//
//	data_inicio=YYYY-MM-DD (obrigatório)
//	data_fim=YYYY-MM-DD    (obrigatório)
//
// Retorna um array de LojaExploradorItem com totalTreinamentos e taxaParticipacao
// calculados para o período informado. Lojas sem nenhuma inscrição no período
// aparecem com zeros.
func ExploradorLojasHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `{"erro": "Método não permitido. Use GET."}`, http.StatusMethodNotAllowed)
		return
	}

	dataInicio := r.URL.Query().Get("data_inicio")
	dataFim := r.URL.Query().Get("data_fim")

	// Valores padrão: ano corrente se não informados
	if dataInicio == "" {
		dataInicio = fmt.Sprintf("%d-01-01", time.Now().Year())
	}
	if dataFim == "" {
		dataFim = fmt.Sprintf("%d-12-31", time.Now().Year())
	}

	items, err := repositories.BuscarExploradorLojas(dataInicio, dataFim)
	if err != nil {
		fmt.Printf("Erro ao buscar explorador de lojas: %v\n", err)
		http.Error(w, `{"erro": "Erro interno ao buscar dados do explorador."}`, http.StatusInternalServerError)
		return
	}

	// Garante array vazio em vez de null quando não há dados
	if items == nil {
		items = make([]models.LojaExploradorItem, 0)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}
