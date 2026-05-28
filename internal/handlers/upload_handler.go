package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"

	//Biblioteca para receber planilhas do excel
	"path/filepath"
	"strings"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
	"github.com/xuri/excelize/v2"
)

//RECEBER A PLANILHA (CSV OU XLSX) DO FRONT E SALVAR AS PRESENÇAS

func UploadPlanilhaHandler(w http.ResponseWriter, r *http.Request) {
	//Configurar o CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	}

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

	file, header, err := r.FormFile("planilha")

	if err != nil {
		http.Error(w, "Erro ao receber o arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	//Descobrir a extensao do arquivo que foi mandado

	extensao := strings.ToLower(filepath.Ext(header.Filename))

	//Variavel para guardar a lista

	var linhas [][]string

	//Verificacao do tipo do formato do arquivo CSV OU XLSX

	switch extensao {
	case ".csv":
		//codigo para ler csv
		leitor := csv.NewReader(file)
		linhas, err = leitor.ReadAll()

		if err != nil {
			http.Error(w, "Erro ao ler as linhas do arquivo CSV.", http.StatusInternalServerError)
			return
		}
	case ".xlsx":
		//codigo para ler xlsx
		f, err := excelize.OpenReader(file)

		if err != nil {
			http.Error(w, "Erro ao ler as linhas do arquivo XLSX.", http.StatusInternalServerError)
		}

		defer f.Close()

		//Pega os dados da tabela do excel
		nomeDaAba := f.GetSheetName(f.GetActiveSheetIndex())
		linhas, err = f.GetRows(nomeDaAba)

		if err != nil {
			http.Error(w, "Erro ao ler as linhas do arquivo XLSX.", http.StatusInternalServerError)
			return
		}
	default:
		//Trava para não aceitar arquivos de outros formatos
		http.Error(w, "Formato inválido. O sistema aceita apenas planilhas .csv ou .xlsx", http.StatusBadRequest)
		return

	}

	// if extensao == ".csv" {
	// 	//codigo para ler csv
	// 	leitor := csv.NewReader(file)
	// 	linhas, err = leitor.ReadAll()

	// 	if err != nil {
	// 		http.Error(w, "Erro ao ler as linhas do arquivo CSV.", http.StatusInternalServerError)
	// 		return
	// 	}
	// } else if extensao == ".xlsx" {
	// 	//codigo para ler xlsx
	// 	f, err := excelize.OpenReader(file)

	// 	if err != nil {
	// 		http.Error(w, "Erro ao ler as linhas do arquivo XLSX.", http.StatusInternalServerError)
	// 	}

	// 	defer f.Close()

	// 	//Pega os dados da tabela do excel
	// 	nomeDaAba := f.GetSheetName(f.GetActiveSheetIndex())
	// 	linhas, err = f.GetRows(nomeDaAba)

	// 	if err != nil {
	// 		http.Error(w, "Erro ao ler as linhas do arquivo XLSX.", http.StatusInternalServerError)
	// 		return
	// 	}
	// } else {
	// 	//Trava para não aceitar arquivos de outros formatos
	// 	http.Error(w, "Formato inválido. O sistema aceita apenas planilhas .csv ou .xlsx", http.StatusBadRequest)
	// 	return
	// }

	//Loop para ler as linhas, tanto para csv ou xlsx

	salvos := 0

	for indice, linha := range linhas {
		//Pula os titulos da tabela que geralmente fica no indice 0
		if indice == 0 {
			continue
			//continue faz a funcao avancar para a proxima linha
		}

		//Trava para pular linhas em branco ou mal formatadas
		if len(linha) < 4 {
			continue
		}

		//Extracao dos dados de cada linha, as colunas e suas informacoes,
		// para isso usamos o slice de linha
		//para dividir ela em pedacos salvamos em variaveis e depois
		//jogando na funcao presenca_repo

		luc := linha[0]
		representante := linha[2]
		//Fazer o status ficar tudo em maisculo para ser aceito no banco
		status := strings.ToUpper(linha[3])

		//Agora salvamos no banco de dados
		err = repositories.SalvarPresencaPlanilha(treinamentoID, luc, representante, status)

		if err != nil {
			fmt.Printf("Aviso: Erro ao salvar presença do LUC %s: %v\n", luc, err)
			continue
		}

		salvos++

	}

	//resposta para o front de sucesso
	mensagem := fmt.Sprintf(`{"mensagem": "Planilha processada com sucesso! %d presenças salvas."}`, salvos)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mensagem))

}
