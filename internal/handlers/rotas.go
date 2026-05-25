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

	// Rota para o check-in automático (Auto-presença via QR Code)
	http.HandleFunc("/api/presencas/confirmar", ConfirmarPresencaHandler)

}
