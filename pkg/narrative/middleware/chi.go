package middleware

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pgquerynarrative/pgquerynarrative/pkg/narrative"
)

// MountChi mounts PgQueryNarrative HTTP handlers on the Chi router under prefix.
// Routes: POST prefix/query/run, POST prefix/report/generate, GET prefix/schema,
// GET prefix/suggestions/queries. Prefix is normalized (no trailing slash).
// Client must not be nil.
func MountChi(r chi.Router, client *narrative.Client, prefix string) {
	prefix = strings.TrimSuffix(prefix, "/")
	if prefix != "" {
		r.Route(prefix, func(r chi.Router) {
			mountChiRoutes(r, client)
		})
	} else {
		mountChiRoutes(r, client)
	}
}

func mountChiRoutes(r chi.Router, client *narrative.Client) {
	r.Post("/query/run", func(w http.ResponseWriter, req *http.Request) {
		HandleRunQuery(client, w, req)
	})
	r.Post("/report/generate", func(w http.ResponseWriter, req *http.Request) {
		HandleGenerateReport(client, w, req)
	})
	r.Get("/schema", func(w http.ResponseWriter, req *http.Request) {
		HandleGetSchema(client, w, req)
	})
	r.Get("/suggestions/queries", func(w http.ResponseWriter, req *http.Request) {
		HandleSuggestionsQueries(client, w, req)
	})
}
