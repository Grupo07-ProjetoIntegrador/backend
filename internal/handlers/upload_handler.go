package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

//RECEBER A PLANILHA DO FRONT E SALVAR AS PRESENÇAS

func UploadPlanilhaHandler(w http.ResponseWriter, r *http.Request) {
	//Configurar o CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	//Trava de segurança para metodo post

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use Post para enviar arquivos.", http.StatusMethodNotAllowed)
		return
	}

	//Pega o id de treinamento na URL
	treinamentoID := r.URL.Query().Get("treinamento_id")

	if treinamentoID == "" {
		http.Error(w, "O ID do treinamento é obrigatório na URL.", http.StatusBadRequest)
		return
	}

	//Limita o tamanho do arquivo
	r.ParseMultipartForm(5 << 20)

	//Recebe o arquivo. Obs: deve se chamar"planilha"

	file, _, err := r.FormFile("planilha")

	if err != nil {
		http.Error(w, "Erro ao receber o arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	//Le o arquivo csv

	leitor := csv.NewReader(file)

	linhas, err := leitor.ReadAll()

	if err != nil {
		http.Error(w, "Erro ao ler as linhas das planilhas", http.StatusInternalServerError)
		return
	}

	//Loop para salvar linha por linha no banco
	salvos := 0

	for indice, linha := range linhas {
		//Pula a primeira linha
		if indice == 0 {
			continue
		}
		//Verificando se a linha tem pelo menos 4 colunas
		if len(linha) < 4 {
			continue
		}

		//Extrai os dados baseados na ordem das colunas da planilha
		luc := linha[0]
		representante := linha[2]
		status := linha[3]

		//Envia os dados para a função do presenca_repo.go

		err = repositories.SalvarPresencaPlanilha(treinamentoID, luc, representante, status)

		if err != nil {
			fmt.Printf("Aviso: Erro ao salvar a presenca do LUC %s: %v", luc, err)
			continue
		}

		//Contagem de sucesso de cada linha Ok
		salvos++
	}

	mensagem := fmt.Sprintf(`{"Mensagem": "Planilha Processada com sucesso! %d presenças salvas."}`, salvos)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mensagem))

}
