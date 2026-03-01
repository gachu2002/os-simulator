package realtime

import (
	"encoding/json"
	"net/http"
)

const maxJSONBodyBytes = 1 << 20

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	respondError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	return false
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	defer func() { _ = r.Body.Close() }()
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid_body", "invalid JSON body")
		return false
	}
	return true
}
