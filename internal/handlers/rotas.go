package handlers

import (
	"net/http"
)

func ConfigurarRotas() {
	//Rota de Cadastro de loja
	http.HandleFunc("/api/lojas/cadastrar", CadastrarLojaHandler)

	// Quando o formulário da tela for submetido, o front-end envia para cá
	http.HandleFunc("/api/treinamentos/cadastrar", CadastrarTreinamentoHandler)

	http.HandleFunc("/api/treinamentos/webhook-forms", ReceberInscricaoForms)

	//Update dos treinamentos
	http.HandleFunc("/api/treinamentos/cadastrar/", UpdateTreinamentosHandler)

	//Rota de listar os treinamentos
	http.HandleFunc("/api/treinamentos", ListarTreinamentosHandler)

	//Rota de deletar treinamento
	http.HandleFunc("/api/treinamentos/deletar", DeletarTreinamentoHandler)

}
