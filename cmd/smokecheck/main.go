package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
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
	if err := getOK(client, *backend+"/lessons"); err != nil {
		fatalf("lessons failed: %v", err)
	}
	attemptID, sessionID, err := startChallenge(client, *backend)
	if err != nil {
		fatalf("challenge start failed: %v", err)
	}
	if err := postOK(client, *backend+"/challenges/grade", map[string]any{"attempt_id": attemptID}); err != nil {
		fatalf("challenge grade failed: %v", err)
	}
	if err := wsSmoke(*backend, sessionID); err != nil {
		fatalf("websocket smoke failed: %v", err)
	}

	if err := getOK(client, *frontend); err != nil {
		fatalf("frontend check failed: %v", err)
	}
}

func startChallenge(client *http.Client, backend string) (string, string, error) {
	b, err := json.Marshal(map[string]any{"lesson_id": "l01-sched-rr-basics", "stage_index": 0})
	if err != nil {
		return "", "", err
	}
	resp, err := client.Post(backend+"/challenges/start", "application/json", bytes.NewReader(b))
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("status=%d body=%s", resp.StatusCode, string(body))
	}
	var out challengeStartResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", err
	}
	if out.AttemptID == "" || out.SessionID == "" {
		return "", "", fmt.Errorf("empty challenge identifiers")
	}
	return out.AttemptID, out.SessionID, nil
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

func wsSmoke(backend, sessionID string) error {
	wsURL, err := toWSURL(backend, sessionID)
	if err != nil {
		return err
	}
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	var connected map[string]any
	if err := conn.ReadJSON(&connected); err != nil {
		return err
	}
	if connected["type"] != "session.snapshot" {
		return fmt.Errorf("unexpected first event type=%v", connected["type"])
	}

	cmd := map[string]any{"type": "command", "command": map[string]any{"name": "step", "count": 1}}
	if err := conn.WriteJSON(cmd); err != nil {
		return err
	}
	var event map[string]any
	if err := conn.ReadJSON(&event); err != nil {
		return err
	}
	if event["type"] != "session.snapshot" {
		return fmt.Errorf("unexpected command event type=%v", event["type"])
	}
	return nil
}

func toWSURL(base, sessionID string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	default:
		return "", fmt.Errorf("unsupported scheme %q", u.Scheme)
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/ws/" + sessionID
	return u.String(), nil
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
