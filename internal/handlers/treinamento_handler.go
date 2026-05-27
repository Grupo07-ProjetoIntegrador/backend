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
	// Liberar o CORS para o Front-end conseguir acessar
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Se o navegador estiver apenas testando a conexão (Preflight OPTIONS), retorna OK
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
}

//Função Get que pega os dados do treinamento e lança na tela de lista

func ListarTreinamentosHandler(w http.ResponseWriter, r *http.Request) {
	// Liberar o CORS para o Front-end conseguir acessar
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Se o navegador estiver apenas testando a conexão (Preflight OPTIONS), retorna OK
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Verificando se esta usando o comando Get
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido. Use Get.", http.StatusMethodNotAllowed)
		return
	}

	//Buscando a lista do banco de dados
	lista, err := repositories.ListarTreinamentos()

	//Verificacao de erro de conexao
	if err != nil {
		http.Error(w, "Erro ao buscar a lista de treinamentos", http.StatusInternalServerError)
		fmt.Println("Erro na listagem:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)

}

func DeletarTreinamentoHandler(w http.ResponseWriter, r *http.Request) {

	//Configuracao do CORS para o front acessar
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Trava que faz a URl aceitar somento o DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Metodo nao permitido. Use DELETE.", http.StatusMethodNotAllowed)

		return

	}

	//Extrai o Id da Url
	id := r.URL.Query().Get("id")

	//Verifica se o Id foi enviado
	if id == "" {
		http.Error(w, "O id de treinamento é obrigatório", http.StatusBadRequest)
		return
	}

	//Chama a função do repositório

	err := repositories.DeletarTreinamento(id)

	//Verifica se tem algum erro do err

	if err != nil {
		http.Error(w, "Erro ao deletar"+err.Error(), http.StatusInternalServerError)
		return
	}

	//Caso de tudo certo vem pra ca
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensagem": "Treinamento deletado com sucesso!"}`))
}
