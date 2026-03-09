// Package services
package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/warmdev17/innogen-backend/internal/models"
)

func RunCode(code string, language string) (*models.ExecuteResponse, error) {
	reqBody := models.ExecuteRequest{
		Language: language,
		Version:  "*",
		Files: []models.File{
			{Content: code},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	pistonURL := os.Getenv("PISTON_URL")
	if pistonURL == "" {
		pistonURL = "http://localhost:2000"
	}

	resp, err := http.Post(
		pistonURL+"/api/v2/piston/execute",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result models.ExecuteResponse
	err = json.NewDecoder(resp.Body).Decode(&result)

	return &result, err
}
