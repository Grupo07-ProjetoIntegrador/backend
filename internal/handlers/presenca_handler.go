package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

type ConfirmarPresencaRequest struct {
	TreinamentoID string  `json:"treinamento_id"`
	Email         string  `json:"email"`
	UserLatitude  float64 `json:"user_latitude"`
	UserLongitude float64 `json:"user_longitude"`
}

func ListarPresencasHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Configuração do CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	treinamentoID := r.URL.Query().Get("treinamento_id")
	if treinamentoID == "" {
		http.Error(w, `{"erro": "ID do treinamento é obrigatório"}`, http.StatusBadRequest)
		return
	}

	presencas, err := repositories.ListarPresencaPorTreinamentos(treinamentoID)

	if err != nil {
		// Se o erro for apenas porque não há linhas, tratamos como sucesso com lista vazia
		// Você pode checar se o erro é 'sql.ErrNoRows' ou se prefere apenas zerar o erro se o seu grupo preferir
		fmt.Printf("🚨 Erro ou aviso ao buscar no banco: %v\n", err)

		// Se quiser que mesmo com erro ele não quebre o front, podemos forçar o envio de uma lista vazia:
		presencas = []models.PresencaResponse{}
	}

	// Se a busca deu certo mas veio nula (sem registros no banco), transformamos em [] para o React não quebrar
	if presencas == nil {
		presencas = []models.PresencaResponse{} // Evita mandar 'null' no JSON, manda '[]'
	}

	// 3. Devolve os dados em formato JSON para o React (Sempre 200 OK se chegou aqui)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(presencas)
}

// CalcularDistanciaHaversine calcula a distância em metros entre dois pontos geográficos
func CalcularDistanciaHaversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000.0 // Raio da Terra em metros
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	rLat1 := lat1 * math.Pi / 180.0
	rLat2 := lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rLat1)*math.Cos(rLat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// ConfirmarPresencaHandler atende o PATCH enviado pelo checkin.html
func ConfirmarPresencaHandler(w http.ResponseWriter, r *http.Request) {
	// CORS Headers - Permite que a página HTML executada no celular faça chamadas para esta API
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Responde requisições OPTIONS pré-voo do CORS
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPatch {
		http.Error(w, "Método não permitido. Use PATCH.", http.StatusMethodNotAllowed)
		return
	}

	var req ConfirmarPresencaRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Erro ao ler os dados enviados", http.StatusBadRequest)
		return
	}

	if req.TreinamentoID == "" || req.Email == "" {
		http.Error(w, "Campos 'treinamento_id' e 'email' são obrigatórios", http.StatusBadRequest)
		return
	}

	// Validação de Geofencing
	var lat, lon float64
	var raio int
	var hasGeofence bool

	err = database.DB.QueryRow(`
		SELECT lt.latitude, lt.longitude, lt.raio_amplitude
		FROM treinamentos t
		INNER JOIN locais_treinamento lt ON t.local_id = lt.id
		WHERE t.id = $1
	`, req.TreinamentoID).Scan(&lat, &lon, &raio)

	if err == nil {
		hasGeofence = true
	} else if err != sql.ErrNoRows {
		fmt.Printf("Aviso ao buscar geofencing no banco: %v\n", err)
	}

	if hasGeofence {
		distancia := CalcularDistanciaHaversine(req.UserLatitude, req.UserLongitude, lat, lon)
		if distancia > float64(raio) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"erro": fmt.Sprintf("Você está fora do perímetro permitido (%d metros). Distância atual: %.2f metros.", raio, distancia),
			})
			return
		}
	}

	nomeParticipante, _ := repositories.BuscarNomeParticipantePorEmail(req.TreinamentoID, req.Email)

	// Tenta atualizar a presença para PRESENTE
	err = repositories.ConfirmarPresencaPorEmail(req.TreinamentoID, req.Email)
	if err != nil {
		// Retorna 404 se não achar a inscrição pendente ou outro erro correspondente
		http.Error(w, fmt.Sprintf("Erro ao confirmar presença: %v", err), http.StatusNotFound)
		return
	}

	if err := notificarPresencaValidada(req.TreinamentoID, req.Email, nomeParticipante); err != nil {
		fmt.Println("Aviso ao enviar e-mail de presença validada:", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Presença confirmada com sucesso para o participante %s!", req.Email)
}

func notificarPresencaValidada(treinamentoID string, email string, nomeParticipante string) error {
	treinamento, err := repositories.BuscarTreinamentoPorID(treinamentoID)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"treinamento_id": treinamentoID,
		"treinamento": map[string]any{
			"id":             treinamento.ID,
			"tema":           treinamento.Tema,
			"descricao":      treinamento.Descricao,
			"objetivo":       treinamento.Objetivo,
			"data":           treinamento.Data,
			"horario_inicio": treinamento.HorarioInicio,
			"horario_fim":    treinamento.HorarioFim,
			"local":          treinamento.Local,
			"segmento_alvo":  treinamento.SegmentoAlvo,
		},
		"destinatario": map[string]any{
			"nome":  nomeParticipante,
			"email": email,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	apiURL := os.Getenv("AUTOMACOES_PUBLIC_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8000"
	}

	resp, err := http.Post(apiURL+"/api/automacoes/notificar-presenca-validada", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("automacao retornou status %d: %s", resp.StatusCode, string(responseBody))
	}

	return nil
}

func CriarPresencaManualHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"erro": "Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}

	var input models.CriarPresencaInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, `{"erro": "JSON inválido"}`, http.StatusBadRequest)
		return
	}

	// Validação dos campos obrigatórios
	if input.TreinamentoID == "" || input.LUC == "" || input.Representante == "" {
		http.Error(w, `{"erro": "Todos os campos são obrigatórios"}`, http.StatusBadRequest)
		return
	}

	input.Status = strings.ToUpper(input.Status)

	// Chama a função do repositório enviando os dados validados
	err = repositories.CriarPresencaManual(input.TreinamentoID, input.LUC, input.Representante, input.Status)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Participante adicionado com sucesso!"})
}

func DeletarPresencaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		return
	}

	// Aceita apenas o método DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, `{"erro": "Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}

	// Pega o ID dos parâmetros da URL (?id=...)
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"erro": "O parâmetro ID é obrigatório"}`, http.StatusMethodNotAllowed)
		return
	}

	// Chama a função do repositório para deletar do banco
	err := repositories.DeletarPresenca(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Presença removida com sucesso!"})
}

// ObterGeofencingTreinamentoHandler retorna as coordenadas e o raio da cerca virtual do treinamento
func ObterGeofencingTreinamentoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

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
		http.Error(w, `{"erro": "ID do treinamento é obrigatório"}`, http.StatusBadRequest)
		return
	}

	var nomeLocal string
	var lat, lon float64
	var raio int

	err := database.DB.QueryRow(`
		SELECT lt.nome_local, lt.latitude, lt.longitude, lt.raio_amplitude
		FROM treinamentos t
		INNER JOIN locais_treinamento lt ON t.local_id = lt.id
		WHERE t.id = $1
	`, treinamentoID).Scan(&nomeLocal, &lat, &lon, &raio)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"erro": "Geofencing não configurado para este treinamento"}`, http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf(`{"erro": "Erro ao buscar geofencing: %v"}`, err), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"nome_local":     nomeLocal,
		"latitude":       lat,
		"longitude":      lon,
		"raio_amplitude": raio,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type CadastrarLocalRequest struct {
	NomeLocal     string  `json:"nome_local"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	RaioAmplitude int     `json:"raio_amplitude"`
}

// CadastrarLocalHandler cadastra um novo ponto de geofencing
func CadastrarLocalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"erro": "Método não permitido. Use POST."}`, http.StatusMethodNotAllowed)
		return
	}

	var req CadastrarLocalRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"erro": "Erro ao ler dados da requisição"}`, http.StatusBadRequest)
		return
	}

	if req.NomeLocal == "" || req.Latitude == 0 || req.Longitude == 0 || req.RaioAmplitude <= 0 {
		http.Error(w, `{"erro": "Todos os campos são obrigatórios e raio deve ser maior que 0"}`, http.StatusBadRequest)
		return
	}

	var id string
	err = database.DB.QueryRow(`
		INSERT INTO locais_treinamento (nome_local, latitude, longitude, raio_amplitude)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, req.NomeLocal, req.Latitude, req.Longitude, req.RaioAmplitude).Scan(&id)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"erro": "Erro ao salvar local no banco: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"mensagem": "Local cadastrado com sucesso!",
		"id":       id,
	})
}

// ListarLocaisHandler retorna todos os locais cadastrados
func ListarLocaisHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `{"erro": "Método não permitido. Use GET."}`, http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, nome_local, latitude, longitude, raio_amplitude
		FROM locais_treinamento
		ORDER BY nome_local ASC
	`)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"erro": "Erro ao buscar locais: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type LocalResponse struct {
		ID            string  `json:"id"`
		NomeLocal     string  `json:"nome_local"`
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		RaioAmplitude int     `json:"raio_amplitude"`
	}

	locais := []LocalResponse{}
	for rows.Next() {
		var loc LocalResponse
		if err := rows.Scan(&loc.ID, &loc.NomeLocal, &loc.Latitude, &loc.Longitude, &loc.RaioAmplitude); err == nil {
			locais = append(locais, loc)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(locais)
}

func EditarPresencaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, `{"erro": "Método não permitido"}`, http.StatusMethodNotAllowed)
		return
	}

	var input models.EditarPresencaInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, `{"erro": "JSON inválido"}`, http.StatusBadRequest)
		return
	}

	if input.ID == "" || input.LUC == "" || input.Representante == "" {
		http.Error(w, `{"erro": "Todos os campos são obrigatórios"}`, http.StatusBadRequest)
		return
	}

	input.Status = strings.ToUpper(input.Status)

	err = repositories.EditarPresenca(input.ID, input.LUC, input.Representante, input.Status)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Participante editado com sucesso!"})
}
