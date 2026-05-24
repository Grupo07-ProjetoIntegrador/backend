package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	// Importando as pastas do seu projeto
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

// CadastrarTreinamentoHandler recebe os dados da tela "Cadastrar Novo Treinamento"
func CadastrarTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Verifica se o Front-end está mandando um POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	// 2. Cria o "molde" vazio usando a Struct que você fez na Fase 1
	var novoTreinamento models.Treinamento

	// 3. Lê o JSON que veio do Front-end e preenche o molde
	err := json.NewDecoder(r.Body).Decode(&novoTreinamento)
	if err != nil {
		http.Error(w, "Erro ao ler os dados do formulário preenchido", http.StatusBadRequest)
		return
	}

	// 4. INTEGRAÇÃO COM O BANCO DE DADOS
	// Chama a função do repositório que salva e devolve o ID gerado
	idGerado, err := repositories.InserirTreinamento(novoTreinamento)

	if err != nil {
		// Se der erro (ex: banco offline), avisa o front-end
		http.Error(w, "Erro ao salvar o treinamento no banco de dados", http.StatusInternalServerError)
		fmt.Println("Erro no repositório:", err)
		return
	}

	// 5. RESPOSTA DE SUCESSO!
	// Retorna Status 201 (Criado) e devolve o UUID para o administrador copiar e usar no Google Forms
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Treinamento '%s' criado com sucesso! O ID para o Google Forms é: %s", novoTreinamento.Tema, idGerado)
}
