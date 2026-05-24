package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ReceberInscricaoForms atende o POST automático vindo do Google Forms
func ReceberInscricaoForms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var inscricao InscricaoFormsRequest
	err := json.NewDecoder(r.Body).Decode(&inscricao)
	if err != nil {
		http.Error(w, "Erro ao ler dados do Forms", http.StatusBadRequest)
		return
	}

	// LÓGICA DE NEGÓCIO A SER IMPLEMENTADA:

	// 1. O código vai buscar no banco de dados se a loja com o 'inscricao.LUC' já existe.
	// Se não existir, ele cria silenciosamente usando o InserirLoja() que fizemos antes.

	// 2. O código vai inserir a presença na tabela 'presencas'.
	// Lembra do nosso SQL? O status_presenca tem o padrão PENDENTE [4].
	// query := `INSERT INTO presencas (treinamento_id, loja_id, nome_participante, status_presenca)
	//           VALUES ($1, $2, $3, 'PENDENTE')`

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Inscrição da loja %s recebida com sucesso. Status: PENDENTE", inscricao.NomeLoja)
}
