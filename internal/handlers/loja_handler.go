package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

// Função que atende a requisição para cadastrar a loja
func CadastrarLojaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//1 - Verifica se o metodo HTTP É POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use o Post.", http.StatusMethodNotAllowed)
		return
	}

	//2- Transforma o JSON recebido na struct Loja Request
	var novaLoja models.Loja

	//Go faz a leitura do arquivo JSON e preencher as variaveis
	err := json.NewDecoder(r.Body).Decode(&novaLoja)

	//Caso de erro essa mensagem aparecera

	if err != nil {
		http.Error(w, "Erro ao ler os dados enviados pelo Front-end", http.StatusBadRequest)
		return
	}

	//Chamando a funcao de inserir loja do arquivo loja_repo.go e usando aqui.
	err = repositories.InserirLoja(novaLoja)

	if err != nil {
		//Erro status 500, avisa ao front que deu erro interno no servidor
		http.Error(w, "Erro ao salvar a loja no banco de dados", http.StatusInternalServerError)
		fmt.Println("Erro no repositorio", err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	fmt.Fprintf(w, "A loja '%s' (LUC: '%s') do segmento '%s' foi processada usando model oficial!", novaLoja.Nome, novaLoja.LUC, novaLoja.Segmento)

}
