// Package judge
package judge

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
	"github.com/warmdev17/innogen-backend/internal/services"
)

func StartWorker() {
	log.Println("Judge Worker started...")

	for {
		result, err := database.Rdb.BLPop(context.Background(), 0*time.Second, "judge_queue").Result()
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

		subIDStr, ok := jobData["submission_id"].(string)
		if !ok {
			log.Println("Invalid submission_id in job payload")
			continue
		}

		go processSubmission(subIDStr)
	}
}

func processSubmission(subID string) {
	var sub models.Submission
	if err := database.DB.First(&sub, "id = ?", subID).Error; err != nil {
		log.Println("Submission not found in DB for ID:", subID)
		return
	}

	database.DB.Model(&sub).Update("status", "Processing")

	var testcases []models.Testcase
	if err := database.DB.Where("problem_id = ?", sub.ProblemID).Find(&testcases).Error; err != nil {
		log.Println("Error finding testcases:", err)
		return
	}

	totalTestcases := len(testcases)
	if totalTestcases == 0 {
		database.DB.Model(&sub).Updates(map[string]any{
			"status":          "Accepted",
			"total_testcases": 0,
			"pass_count":      0,
		})
		log.Printf("Done submission %s: Accepted (No Testcases)\n", subID)
		return
	}

	passCount := 0
	status := "Accepted"
	errorMsg := ""
	maxRuntime := 0

	for _, tc := range testcases {
		resp, err := services.RunCodeWithInput(sub.Code, sub.Language, tc.InputData)
		if err != nil {
			status = "System Error"
			errorMsg = err.Error()
			break
		}

		if resp.Run.Code != 0 {
			status = "Runtime Error"
			if resp.Run.Stderr != "" {
				errorMsg = resp.Run.Stderr
			} else {
				errorMsg = resp.Run.Stdout
			}
			break
		}

		output := strings.TrimSpace(resp.Run.Stdout)
		expected := strings.TrimSpace(tc.ExpectedOutput)

		if output != expected {
			status = "Wrong Answer"
			errorMsg = "Expected: " + expected + "\nGot: " + output
			break
		}

		passCount++
	}

	database.DB.Model(&sub).Updates(map[string]any{
		"status":          status,
		"error_message":   errorMsg,
		"pass_count":      passCount,
		"total_testcases": totalTestcases,
		"runtime_ms":      maxRuntime,
	})

	log.Printf("Done submission %s: %s (Passed %d/%d)\n", subID, status, passCount, totalTestcases)
}
