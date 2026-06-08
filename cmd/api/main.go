package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/handlers"
	"github.com/joho/godotenv"
)

// 1. A FUNÇÃO GLOBAL DE CORS FICA AQUI (Fora da main)
func aplicarCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permite requisições de qualquer origem (como o localhost:5173 do seu React)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Trata a requisição de "pré-voo" (OPTIONS) que o navegador faz antes do POST/PATCH/DELETE
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Passa a requisição adiante para os seus handlers normais
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	fmt.Println("Backend rodando e .env carregado!")

	database.ConectarSupabase()

	if database.DB != nil {
		defer database.DB.Close()
	}

	handlers.ConfigurarRotas()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API do Módulo de Treinamentos rodando com sucesso!")
	})

	fmt.Println("Servidor rodando na porta: http://localhost:8080")

	// 2. MODIFICAÇÃO AQUI: Envolvemos o roteador padrão (http.DefaultServeMux) no nosso Middleware de CORS
	err := http.ListenAndServe(":8080", aplicarCORS(http.DefaultServeMux))

	if err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
