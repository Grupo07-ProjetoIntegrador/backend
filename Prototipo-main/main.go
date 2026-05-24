package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Rota raiz - redireciona para lista de treinamentos
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/treinamentos/lista", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	// ========== ROTAS DE TREINAMENTOS ==========

	// Lista de treinamentos
	http.HandleFunc("/treinamentos/lista", serveHTML("src/imports/lista.html"))

	// Calendário de treinamentos
	http.HandleFunc("/treinamentos/calendario", serveHTML("src/imports/calendario.html"))

	// Formulário: novo treinamento
	http.HandleFunc("/treinamentos/novo", serveHTML("src/imports/formulariotreinamento.html"))

	// Formulário: editar treinamento
	http.HandleFunc("/treinamentos/editar", serveHTML("src/imports/formulariotreinamento.html"))

	// Detalhes de um treinamento específico
	http.HandleFunc("/treinamentos/detalhes", serveHTML("src/imports/detalhesdotreinamento.html"))

	// ========== ROTA DE ENGAJAMENTO ==========

	// Dashboard de engajamento das lojas
	http.HandleFunc("/engajamento", serveHTML("src/imports/engajamento.html"))

	// ========== ROTAS DE LOJAS ==========

	// Detalhes de uma loja específica
	http.HandleFunc("/lojas/detalhes", serveHTML("src/imports/detalhesdaloja.html"))

	// ========== ARQUIVOS ESTÁTICOS ==========

	// Servir arquivos estáticos (CSS, JS, imagens)
	// Certifique-se de ter uma pasta "static" com os recursos necessários
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Servir assets do build (imagens, logos)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	// ========== INICIAR SERVIDOR ==========

	port := ":8080"
	log.Printf("🚀 Servidor rodando em http://localhost%s", port)
	log.Printf("📍 Página inicial: http://localhost%s/treinamentos/lista", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// serveHTML retorna um handler que serve um arquivo HTML específico
func serveHTML(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define o tipo de conteúdo como HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Resolve o caminho absoluto do arquivo
		absPath, err := filepath.Abs(filename)
		if err != nil {
			log.Printf("❌ Erro ao resolver caminho do arquivo %s: %v", filename, err)
			http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
			return
		}

		// Serve o arquivo HTML
		http.ServeFile(w, r, absPath)
	}
}
