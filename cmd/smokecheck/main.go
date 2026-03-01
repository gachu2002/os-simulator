package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type challengeStartResponse struct {
	AttemptID string `json:"attempt_id"`
	SessionID string `json:"session_id"`
}

func main() {
	backend := flag.String("backend", os.Getenv("BACKEND_BASE_URL"), "backend base URL")
	frontend := flag.String("frontend", os.Getenv("FRONTEND_BASE_URL"), "frontend base URL")
	flag.Parse()

	if *backend == "" || *frontend == "" {
		fmt.Fprintln(os.Stderr, "backend and frontend URLs are required")
		os.Exit(1)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	if err := getOK(client, *backend+"/healthz"); err != nil {
		fatalf("healthz failed: %v", err)
	}
	if err := getOK(client, *backend+"/curriculum/v3"); err != nil {
		fatalf("curriculum/v3 failed: %v", err)
	}
	attemptID, err := startChallengeV3(client, *backend)
	if err != nil {
		fatalf("challenge start v3 failed: %v", err)
	}
	if err := postOK(client, *backend+"/challenges/action/v3", map[string]any{"attempt_id": attemptID, "action": "execute_instruction", "count": 1}); err != nil {
		fatalf("challenge action v3 failed: %v", err)
	}
	if err := postOK(client, *backend+"/challenges/submit/v3", map[string]any{"attempt_id": attemptID}); err != nil {
		fatalf("challenge submit v3 failed: %v", err)
	}

	if err := getOK(client, *frontend); err != nil {
		fatalf("frontend check failed: %v", err)
	}
}

func startChallengeV3(client *http.Client, backend string) (string, error) {
	b, err := json.Marshal(map[string]any{"lesson_id": "l01-process-basics"})
	if err != nil {
		return "", err
	}
	resp, err := client.Post(backend+"/challenges/start/v3", "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status=%d body=%s", resp.StatusCode, string(body))
	}
	var out challengeStartResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.AttemptID == "" {
		return "", fmt.Errorf("empty challenge identifiers")
	}
	return out.AttemptID, nil
}

func getOK(client *http.Client, endpoint string) error {
	resp, err := client.Get(endpoint)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status=%d body=%s", resp.StatusCode, string(b))
	}
	return nil
}

func postOK(client *http.Client, endpoint string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := client.Post(endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status=%d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
