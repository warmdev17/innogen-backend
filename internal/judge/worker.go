// Package judge
package judge

import (
	"encoding/json"
	"log"
	"time"

	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"github.com/warmdev17/innogen-backend/internal/services"
)

func StartWorker() {
	log.Println("Judge Worker started...")

	for {
		result, err := database.Rdb.BLPop(database.Ctx, 0*time.Second, "judge_queue").Result()
		if err != nil {
			log.Println("Redis error:", err)
			continue
		}

		payload := result[1]
		log.Println("New job received:", payload)

		var jobData map[string]any
		if err := json.Unmarshal([]byte(payload), &jobData); err != nil {
			log.Println("Invalid JSON:", err)
			continue
		}

		subID := uint(jobData["submission_id"].(float64))

		processSubmission(subID)
	}
}

func processSubmission(subID uint) {
	var sub models.Submission
	if err := database.DB.First(&sub, subID).Error; err != nil {
		log.Println("Submission not found:", subID)
		return
	}

	database.DB.Model(&sub).Update("status", "processing")

	resp, err := services.RunCode(sub.Code, sub.Language)
	status := "accepted"
	output := ""

	if err != nil {
		status = "system_error"
		output = err.Error()
	} else if resp.Run.Code != 0 {
		status = "runtime_error"
		if resp.Run.Stderr != "" {
			output = resp.Run.Stderr
		} else {
			output = resp.Run.Stdout
		}
	} else {
		output = resp.Run.Stdout
	}

	database.DB.Model(&sub).Updates(map[string]any{
		"status": status,
	})

	log.Printf("Done submission %d: %s\nOutput: %s", subID, status, output)
}
