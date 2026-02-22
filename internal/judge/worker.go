// Package judge
package judge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/warmdev17/innogen-backend/internal/database"
	"github.com/warmdev17/innogen-backend/internal/models"
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

	output, err := runCodeInDocker(sub.Code, sub.Language)
	status := "accepted"
	if err != nil {
		status = "runtime_error"
		output = err.Error()
	}

	database.DB.Model(&sub).Updates(map[string]any{
		"status": status,
	})

	log.Printf("Done submission %d: %s\nOutput: %s", subID, status, output)
}

func runCodeInDocker(code string, lang string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer func() {
		if err := cli.Close(); err != nil {
			log.Println("Error closing logs:", err)
		}
	}()

	imageName := "python:3.10-alpine"
	cmd := []string{"python", "-c", code}

	if lang == "go" {
		return "", fmt.Errorf("go runner not implemented yet, try python")
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   cmd,
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", err
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Println("Error closing logs:", err)
		}
	}()

	var outBuf, errBuf bytes.Buffer

	_, err = stdcopy.StdCopy(&outBuf, &errBuf, out)
	if err != nil {
		return "", err
	}

	finalOutput := outBuf.String()
	if errBuf.Len() > 0 {
		finalOutput += "\n--- Stderr ---\n" + errBuf.String()
	}
	// ----------------

	_ = cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})

	return finalOutput, nil
}
