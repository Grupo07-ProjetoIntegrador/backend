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

	// Rota de listar os treinamentos
	http.HandleFunc("/api/treinamentos", ListarTreinamentosHandler)

	// Rota para gerar formulario manualmente
	http.HandleFunc("/api/treinamentos/gerar-formulario", GerarFormularioTreinamentoHandler)

	// Rota para buscar link do formulario
	http.HandleFunc("/api/treinamentos/formulario", BuscarFormularioTreinamentoHandler)

	// Rota para apagar link do formulario
	http.HandleFunc("/api/treinamentos/apagar-formulario", ApagarFormularioTreinamentoHandler)

	// Rota para regerar formulario
	http.HandleFunc("/api/treinamentos/regerar-formulario", RegerarFormularioTreinamentoHandler)

	// Rota de deletar treinamento
	http.HandleFunc("/api/treinamentos/deletar", DeletarTreinamentoHandler)

	// Rotas OAuth Google
	http.HandleFunc("/api/oauth/google/start", GoogleOAuthStartHandler)
	http.HandleFunc("/api/oauth/google/callback", GoogleOAuthCallbackHandler)
	http.HandleFunc("/api/oauth/google/status", GoogleOAuthStatusHandler)
	http.HandleFunc("/api/oauth/google/disconnect", GoogleOAuthDisconnectHandler)
}
