package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type apiError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func respondError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	reqID, _ := requestIDFromContext(r.Context())
	respondJSON(w, status, apiError{Code: code, Message: message, RequestID: reqID})
}

func withCORS(allowed map[string]struct{}, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if isOriginAllowed(allowed, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			if origin != "" && !isOriginAllowed(allowed, origin) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type ctxKey string

const requestIDKey ctxKey = "request_id"

func withRequestID(next http.Handler) http.Handler {
	var seq atomic.Uint64
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if reqID == "" {
			reqID = fmt.Sprintf("req-%08d", seq.Add(1))
		}
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestIDKey).(string)
	return v, ok
}

func allowedOriginsFromEnv(raw string) map[string]struct{} {
	out := map[string]struct{}{}
	if strings.TrimSpace(raw) == "" {
		out["http://localhost:5173"] = struct{}{}
		out["http://127.0.0.1:5173"] = struct{}{}
		out["http://localhost:4173"] = struct{}{}
		out["http://127.0.0.1:4173"] = struct{}{}
		out["https://localhost:5173"] = struct{}{}
		out["https://127.0.0.1:5173"] = struct{}{}
		out["https://localhost:4173"] = struct{}{}
		out["https://127.0.0.1:4173"] = struct{}{}
		return out
	}
	for _, part := range strings.Split(raw, ",") {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		out[origin] = struct{}{}
	}
	return out
}

func isOriginAllowed(allowed map[string]struct{}, origin string) bool {
	if origin == "" {
		return true
	}
	_, ok := allowed[origin]
	return ok
}
