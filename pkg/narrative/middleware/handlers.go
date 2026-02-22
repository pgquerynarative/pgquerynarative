package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

	suggestions "github.com/pgquerynarrative/pgquerynarrative/api/gen/suggestions"
	"github.com/pgquerynarrative/pgquerynarrative/pkg/narrative"
)

// writeJSON sets Content-Type and encodes v. It does not write status code.
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error body and status code.
func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// HandleRunQuery handles POST body {"sql":"...", "limit": N}. Uses default limit if <= 0.
func HandleRunQuery(client *narrative.Client, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var body struct {
		SQL   string `json:"sql"`
		Limit int    `json:"limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	limit := body.Limit
	if limit <= 0 {
		limit = narrative.DefaultRunQueryLimit
	}
	result, err := client.RunQuery(r.Context(), body.SQL, limit)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, result)
}

// HandleGenerateReport handles POST body {"sql":"..."}.
func HandleGenerateReport(client *narrative.Client, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var body struct {
		SQL string `json:"sql"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	report, err := client.GenerateReport(r.Context(), body.SQL)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, report)
}

// HandleGetSchema handles GET and returns schema result.
func HandleGetSchema(client *narrative.Client, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	result, err := client.GetSchema(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, result)
}

// HandleSuggestionsQueries handles GET with query params intent and limit.
func HandleSuggestionsQueries(client *narrative.Client, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	payload := &suggestions.QueriesPayload{Limit: 5}
	if s := r.URL.Query().Get("intent"); s != "" {
		payload.Intent = &s
	}
	if s := r.URL.Query().Get("limit"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil && n > 0 {
			payload.Limit = int32(n)
		}
	}
	result, err := client.SuggestionsService().Queries(r.Context(), payload)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, result)
}
