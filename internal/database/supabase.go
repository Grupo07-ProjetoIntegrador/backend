package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConectarSupabase() {

	// 1. Carrega as variáveis do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. Usando variáveis de ambiente do sistema.")
	}

	// 2. Pega a URL do Supabase
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("A variável DATABASE_URL não foi definida.")
	}

	// 3. Abre a conexão com o PostgreSQL do Supabase
	banco, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Erro ao abrir a conexão com o Supabase:", err)
	}

	// 4. Testa se a conexão realmente funcionou (Ping)
	err = banco.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar (ping) no Supabase:", err)
	}

	fmt.Println("Conexão com o Supabase estabelecida com sucesso!")

	// Salva a conexão na nossa variável global
	DB = banco
}
