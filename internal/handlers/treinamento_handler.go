package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	// Importando as pastas do seu projeto
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

// ListarTreinamentosHandler retorna todos os treinamentos em JSON
func ListarTreinamentosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	treinamentos, err := repositories.ListarTreinamentos()
	if err != nil {
		fmt.Println("Erro ao listar treinamentos:", err)
		http.Error(w, "Erro ao buscar treinamentos no banco de dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(treinamentos)
}

// CadastrarTreinamentoHandler recebe os dados da tela "Cadastrar Novo Treinamento"
func CadastrarTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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

	// 6. DISPARA A GERAÇÃO DO GOOGLE FORMS E ENVIO DE E-MAIL EM SEGUNDO PLANO
	go func(id, tema string) {
		payload := map[string]string{
			"treinamento_id": id,
			"tema":           tema,
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("[Automação] Erro ao serializar JSON para gerar forms: %v\n", err)
			return
		}

		apiURL := "http://localhost:8000/api/automacoes/gerar-forms"
		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Printf("[Automação] Erro ao chamar endpoint de gerar forms: %v\n", err)
			return
		}
		defer resp.Body.Close()
		fmt.Printf("[Automação] Resposta do script de gerar forms: %s\n", resp.Status)
	}(idGerado, novoTreinamento.Tema)
}
