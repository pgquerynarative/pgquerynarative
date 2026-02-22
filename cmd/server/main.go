// Package main provides the entry point for the PgQueryNarrative server.
// It initializes the HTTP server, sets up database connections, and starts
// serving API and web UI requests.
package main

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pgquerynarrative/pgquerynarrative/api/gen/http/queries/server"
	reportsServer "github.com/pgquerynarrative/pgquerynarrative/api/gen/http/reports/server"
	schemaServer "github.com/pgquerynarrative/pgquerynarrative/api/gen/http/schema/server"
	suggestionsServer "github.com/pgquerynarrative/pgquerynarrative/api/gen/http/suggestions/server"
	"github.com/pgquerynarrative/pgquerynarrative/api/gen/queries"
	"github.com/pgquerynarrative/pgquerynarrative/api/gen/reports"
	schema "github.com/pgquerynarrative/pgquerynarrative/api/gen/schema"
	suggestions "github.com/pgquerynarrative/pgquerynarrative/api/gen/suggestions"
	"github.com/pgquerynarrative/pgquerynarrative/app/config"
	"github.com/pgquerynarrative/pgquerynarrative/pkg/narrative"
	"github.com/pgquerynarrative/pgquerynarrative/web"
	goahttp "goa.design/goa/v3/http"
)

const gracefulTimeout = 10 * time.Second

// main is the application entry point. It loads config, creates the narrative
// client (which owns DB pools, runner, LLM, and services), wires Goa endpoints
// and web UI to that client, and runs the HTTP server with graceful shutdown.
func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client, err := narrative.NewClient(ctx, narrative.FromAppConfig(cfg))
	if err != nil {
		log.Fatalf("failed to create narrative client: %v", err)
	}
	defer client.Close()

	logger := log.New(os.Stdout, "[pgquerynarrative] ", log.LstdFlags)

	queriesEndpoints := queries.NewEndpoints(client.QueriesService())
	reportsEndpoints := reports.NewEndpoints(client.ReportsService())
	schemaEndpoints := schema.NewEndpoints(client.SchemaService())
	suggestionsEndpoints := suggestions.NewEndpoints(client.SuggestionsService())

	// Configure HTTP server
	httpServer := setupHTTPServer(cfg, queriesEndpoints, reportsEndpoints, schemaEndpoints, suggestionsEndpoints, logger)

	// Start server in a goroutine
	go func() {
		logger.Printf("🚀 Server listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Println("🛑 Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Printf("shutdown error: %v", err)
	} else {
		logger.Println("✅ Server stopped gracefully")
	}
}

// setupHTTPServer configures and returns an HTTP server with:
// - API routes (via Goa framework) at /api/v1/*
// - Web UI routes at /, /query, /saved, /reports
// - Static file serving at /static/*
func setupHTTPServer(
	cfg config.Config,
	queriesEndpoints *queries.Endpoints,
	reportsEndpoints *reports.Endpoints,
	schemaEndpoints *schema.Endpoints,
	suggestionsEndpoints *suggestions.Endpoints,
	logger *log.Logger,
) *http.Server {
	// Create Goa HTTP muxer for API routes
	mux := goahttp.NewMuxer()
	dec := goahttp.RequestDecoder
	enc := goahttp.ResponseEncoder
	errHandler := func(ctx context.Context, w http.ResponseWriter, err error) {
		_ = goahttp.ErrorEncoder(enc, nil)(ctx, w, err)
	}

	// Mount queries API endpoints
	queriesHTTP := server.New(queriesEndpoints, mux, dec, enc, errHandler, nil)
	server.Mount(mux, queriesHTTP)

	// Mount reports API endpoints
	reportsHTTP := reportsServer.New(reportsEndpoints, mux, dec, enc, errHandler, nil)
	reportsServer.Mount(mux, reportsHTTP)

	// Mount schema API endpoints
	schemaHTTP := schemaServer.New(schemaEndpoints, mux, dec, enc, errHandler, nil)
	schemaServer.Mount(mux, schemaHTTP)

	// Mount suggestions API endpoints
	suggestionsHTTP := suggestionsServer.New(suggestionsEndpoints, mux, dec, enc, errHandler, nil)
	suggestionsServer.Mount(mux, suggestionsHTTP)

	// Create web UI handlers
	webHandlers := web.NewHandlers(queriesEndpoints, reportsEndpoints)

	// Create standard HTTP mux for web routes
	// (Goa mux doesn't support HandleFunc, so we use standard mux for web UI)
	webMux := http.NewServeMux()

	// Serve static files (CSS, JS)
	fs := http.FileServer(http.Dir("./web/static"))
	webMux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Web UI routes
	webMux.HandleFunc("/", webHandlers.Home)
	webMux.HandleFunc("/query", webHandlers.QueryPage)
	webMux.HandleFunc("/saved", webHandlers.SavedQueries)
	webMux.HandleFunc("/reports", webHandlers.Reports)

	// Web API endpoints (form handlers that call Goa endpoints)
	webMux.HandleFunc("/web/query/run", webHandlers.RunQuery)
	webMux.HandleFunc("/web/reports/generate", webHandlers.GenerateReport)

	// Combine API and web routes
	// API routes (/api/*) handled by Goa muxer
	// Web routes (everything else) handled by standard mux
	combinedMux := http.NewServeMux()
	combinedMux.Handle("/api/", mux)
	combinedMux.Handle("/", webMux)

	// Wrap with request logging middleware
	loggedHandler := requestLoggingMiddleware(combinedMux, logger)

	// Create and configure HTTP server
	return &http.Server{
		Addr:         cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port),
		Handler:      loggedHandler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
}

// requestLoggingMiddleware logs each HTTP request: method, path, client IP, status, duration.
// For 4xx/5xx responses it also logs the response body (error message) so the console shows why the request failed.
func requestLoggingMiddleware(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		clientIP := clientIPFromRequest(r)
		path := r.URL.Path
		if path == "" {
			path = "/"
		}
		method := r.Method

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK, logger: logger, method: method, path: path}
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		logger.Printf("%s %s %s %d %s", clientIP, method, path, wrapped.statusCode, duration.Round(time.Millisecond))
		wrapped.logErrorIfAny()
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
	capture    bool
	logger     *log.Logger
	method     string
	path       string
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	if code >= 400 {
		rw.capture = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(p []byte) (n int, err error) {
	if rw.capture && rw.body.Len() < 2048 {
		rw.body.Write(p)
	}
	return rw.ResponseWriter.Write(p)
}

func (rw *responseWriter) logErrorIfAny() {
	if rw.statusCode < 400 || rw.logger == nil {
		return
	}
	body := strings.TrimSpace(rw.body.String())
	if body == "" {
		rw.logger.Printf("error response %d %s %s", rw.statusCode, rw.method, rw.path)
		return
	}
	const max = 512
	if len(body) > max {
		body = body[:max] + "..."
	}
	body = strings.ReplaceAll(body, "\n", " ")
	rw.logger.Printf("error response %d %s %s | %s", rw.statusCode, rw.method, rw.path, body)
}

func clientIPFromRequest(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
