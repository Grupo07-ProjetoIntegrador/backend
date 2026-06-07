package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
)

// GerarDossieLojaHandler queries the DB and forwards the PDF generation request to Python
func GerarDossieLojaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	lojaID := r.URL.Query().Get("loja_id")
	dataInicio := r.URL.Query().Get("data_inicio")
	dataFim := r.URL.Query().Get("data_fim")

	if lojaID == "" || dataInicio == "" || dataFim == "" {
		http.Error(w, "Parâmetros 'loja_id', 'data_inicio' e 'data_fim' são obrigatórios.", http.StatusBadRequest)
		return
	}

	// 1. Fetch store info
	var storeName, storeLuc, storeSegment string
	err := database.DB.QueryRow(`
		SELECT nome, luc, segmento 
		FROM lojas 
		WHERE id = $1
	`, lojaID).Scan(&storeName, &storeLuc, &storeSegment)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar loja: %v", err), http.StatusNotFound)
		return
	}

	dadosLoja := map[string]string{
		"nome":     storeName,
		"luc":      storeLuc,
		"segmento": storeSegment,
		"gerente":  "Responsável da Loja",
	}

	// 2. Fetch trainings and registrations in the period
	rows, err := database.DB.Query(`
		SELECT 
			t.id::text, 
			t.tema, 
			TO_CHAR(t.horario_inicio, 'DD/MM/YYYY') as data_formatada,
			p.nome_participante, 
			p.status_presenca
		FROM presencas p
		INNER JOIN treinamentos t ON p.treinamento_id = t.id
		WHERE p.loja_id = $1 AND t.data BETWEEN $2 AND $3
		ORDER BY t.horario_inicio DESC
	`, lojaID, dataInicio, dataFim)

	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao consultar histórico de treinamentos: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Organize row data grouped by Training ID
	type TempTraining struct {
		Tema      string
		Data      string
		Presentes []string
		Ausentes  []string
	}
	trainingsMap := make(map[string]*TempTraining)
	var trainingsOrder []string

	for rows.Next() {
		var tID, tema, dataFormatada, nomePart, status string
		if err := rows.Scan(&tID, &tema, &dataFormatada, &nomePart, &status); err != nil {
			continue
		}

		if _, exists := trainingsMap[tID]; !exists {
			trainingsMap[tID] = &TempTraining{
				Tema:      tema,
				Data:      dataFormatada,
				Presentes: []string{},
				Ausentes:  []string{},
			}
			trainingsOrder = append(trainingsOrder, tID)
		}

		if status == "PRESENTE" {
			trainingsMap[tID].Presentes = append(trainingsMap[tID].Presentes, nomePart)
		} else {
			trainingsMap[tID].Ausentes = append(trainingsMap[tID].Ausentes, nomePart)
		}
	}

	historicoList := []map[string]any{}
	for _, tID := range trainingsOrder {
		t := trainingsMap[tID]
		historicoList = append(historicoList, map[string]any{
			"tema":      t.Tema,
			"data":      t.Data,
			"presentes": t.Presentes,
			"ausentes":  t.Ausentes,
		})
	}

	// 3. Prepare payload for Python service
	payload := map[string]any{
		"dados_loja":             dadosLoja,
		"period": map[string]string{
			"data_inicio": dataInicio,
			"data_fim":    dataFim,
		},
		"historico_treinamentos": historicoList,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Erro ao serializar payload.", http.StatusInternalServerError)
		return
	}

	// 4. Request PDF from Python service with safe 35-second timeout
	client := &http.Client{
		Timeout: 35 * time.Second,
	}
	pythonURL := automacoesBaseURL() + "/api/automacoes/pdf/dossie"
	resp, err := client.Post(pythonURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao conectar ao serviço de PDF: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Serviço de PDF retornou erro: %s", string(body)), http.StatusBadGateway)
		return
	}

	// 5. Return PDF bytes
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=dossie_%s.pdf", storeLuc))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}

// GerarChamadaTreinamentoHandler queries DB and forwards PDF generation request to Python
func GerarChamadaTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	treinamentoID := r.URL.Query().Get("treinamento_id")
	if treinamentoID == "" {
		http.Error(w, "Parâmetro 'treinamento_id' é obrigatório.", http.StatusBadRequest)
		return
	}

	// 1. Fetch training details
	var tema, dataFormatada, horaFormatada, local, segmentoAlvo, instrutor string
	var capacidadeMax int
	err := database.DB.QueryRow(`
		SELECT 
			tema, 
			TO_CHAR(horario_inicio, 'DD/MM/YYYY'), 
			TO_CHAR(horario_inicio, 'HH24:MI'), 
			local, 
			capacidade_maxima, 
			segmento_alvo,
			responsavel
		FROM treinamentos 
		WHERE id = $1
	`, treinamentoID).Scan(&tema, &dataFormatada, &horaFormatada, &local, &capacidadeMax, &segmentoAlvo, &instrutor)

	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar treinamento: %v", err), http.StatusNotFound)
		return
	}

	dadosTreinamento := map[string]any{
		"id":            treinamentoID,
		"tema":          tema,
		"data":          dataFormatada,
		"hora":          horaFormatada,
		"local":         local,
		"carga_horaria": "2 horas", // Default value
		"segmento_alvo": segmentoAlvo,
		"instrutor":     instrutor,
	}

	// 2. Fetch presence records
	rows, err := database.DB.Query(`
		SELECT 
			p.nome_participante, 
			COALESCE(p.cargo, 'Representante'), 
			l.nome, 
			p.status_presenca,
			COALESCE(TO_CHAR(p.data_registro, 'HH24:MI'), '--:--')
		FROM presencas p
		INNER JOIN lojas l ON p.loja_id = l.id
		WHERE p.treinamento_id = $1
		ORDER BY p.nome_participante ASC
	`, treinamentoID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar presença: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	presentesList := []map[string]string{}
	ausentesList := []map[string]string{}

	for rows.Next() {
		var nome, cargo, loja, status, horario string
		if err := rows.Scan(&nome, &cargo, &loja, &status, &horario); err != nil {
			continue
		}

		person := map[string]string{
			"nome":    nome,
			"cargo":   cargo,
			"loja":    loja,
			"horario": horario,
		}

		if status == "PRESENTE" {
			presentesList = append(presentesList, person)
		} else {
			ausentesList = append(ausentesList, person)
		}
	}

	// 3. Prepare payload for Python service
	payload := map[string]any{
		"dados_treinamento": dadosTreinamento,
		"presentes":         presentesList,
		"ausentes":          ausentesList,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Erro ao serializar payload.", http.StatusInternalServerError)
		return
	}

	// 4. Request PDF from Python service with safe 35-second timeout
	client := &http.Client{
		Timeout: 35 * time.Second,
	}
	pythonURL := automacoesBaseURL() + "/api/automacoes/pdf/chamada"
	resp, err := client.Post(pythonURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao conectar ao serviço de PDF: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Serviço de PDF retornou erro: %s", string(body)), http.StatusBadGateway)
		return
	}

	// 5. Return PDF bytes
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=ata_%s.pdf", treinamentoID))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}
