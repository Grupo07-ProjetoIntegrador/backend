package models

type InscricaoFormsRequest struct {
	TreinamentoID     string `json:"treinamento_id"`     // ID do treinamento que está acontecendo
	LUC               string `json:"luc"`                // LUC digitado no Forms
	NomeLoja          string `json:"nome_loja"`          // Nome da Loja digitado
	NomeRepresentante string `json:"nome_representante"` // Representante digitado
}
