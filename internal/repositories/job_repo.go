package repositories

import (
	"encoding/json"
	"github.com/Grupo07-ProjetoIntegrador/backend/internal/database"
)

func EnfileirarJob(taskType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(`
		INSERT INTO job_queue (task_type, payload, status)
		VALUES ($1, $2, 'pending')
	`, taskType, string(payloadBytes))
	
	return err
}
