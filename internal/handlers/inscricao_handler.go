package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

// ReceberInscricaoForms atende o POST automático vindo do Google Forms
func ReceberInscricaoForms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var inscricao models.InscricaoFormsRequest
	err := json.NewDecoder(r.Body).Decode(&inscricao)
	if err != nil {
		http.Error(w, "Erro ao ler dados do Forms", http.StatusBadRequest)
		return
	}

	// 1. O código busca no banco de dados a loja ativa pelo nome selecionado no Forms.
	lojaID, err := repositories.BuscarLojaAtivaPorNome(inscricao.NomeLoja)
	if err != nil {
		http.Error(w, fmt.Sprintf("Loja invalida: %v", err), http.StatusBadRequest)
		return
	}

	// 2. Insere a presença na tabela 'presencas' com o status PENDENTE, incluindo as novas colunas
	err = repositories.InserirPresencaPendente(
		inscricao.TreinamentoID,
		lojaID,
		inscricao.NomeRepresentante,
		inscricao.Email,
		inscricao.Telefone,
		inscricao.Cargo,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao registrar presença pendente: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Inscrição de '%s' (Loja: %s) recebida com sucesso. Status: PENDENTE", inscricao.NomeRepresentante, inscricao.NomeLoja)
}
