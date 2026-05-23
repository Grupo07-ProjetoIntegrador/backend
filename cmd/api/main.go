package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"github.com/joho/godotenv" // Importando uma biblioteca externa
)

func main() {
	godotenv.Load()
	fmt.Println("Backend rodando e .env carregado!")

	database.ConectarSupabase()

	if database.DB != nil {
		defer database.DB.Close()
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API do Módulo de Treinamentos rodando com sucesso!")
	})

	fmt.Println("Servidor rodando na porta: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
