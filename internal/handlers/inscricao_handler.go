package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Grupo07-ProjetoIntegrador/backend/internal/models"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/repositories"
)

// ReceberInscricaoForms atende o POST automático vindo do Google Forms
func ReceberInscricaoForms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var inscricao models.InscricaoFormsRequest
	err := json.NewDecoder(r.Body).Decode(&inscricao)
	if err != nil {
		http.Error(w, "Erro ao ler dados do Forms", http.StatusBadRequest)
		return
	}

	// 1. O código busca no banco de dados a loja ativa pelo nome selecionado no Forms.
	lojaID, err := repositories.BuscarLojaAtivaPorNome(inscricao.NomeLoja)
	if err != nil {
		// Se a loja não for encontrada ou inativa, faz fallback para a loja genérica "Outra Loja (Não listada)"
		var errFallback error
		lojaID, errFallback = repositories.BuscarOuCriarLoja("9999", "Outra Loja (Não listada)")
		if errFallback != nil {
			http.Error(w, fmt.Sprintf("Erro ao resolver loja de fallback: %v", errFallback), http.StatusInternalServerError)
			return
		}
	}

	// 2. Insere a presença na tabela 'presencas' com o status PENDENTE, incluindo as novas colunas
	err = repositories.InserirPresencaPendente(
		inscricao.TreinamentoID,
		lojaID,
		inscricao.NomeRepresentante,
		inscricao.Email,
		inscricao.Telefone,
		inscricao.Cargo,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao registrar presença pendente: %v", err), http.StatusInternalServerError)
		return
	}

	// Dispara o e-mail de confirmação via serviço de automação Python
	if err := notificarInscricaoConfirmada(inscricao.TreinamentoID, inscricao.Email, inscricao.NomeRepresentante); err != nil {
		fmt.Println("Aviso ao enviar e-mail de confirmação de inscrição:", err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Inscrição de '%s' (Loja: %s) recebida com sucesso. Status: PENDENTE", inscricao.NomeRepresentante, inscricao.NomeLoja)
}

func notificarInscricaoConfirmada(treinamentoID string, email string, nomeParticipante string) error {
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

	resp, err := http.Post(apiURL+"/api/automacoes/notificar-inscricao-confirmada", "application/json", bytes.NewBuffer(body))
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
