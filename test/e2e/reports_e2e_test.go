package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	goahttp "goa.design/goa/v3/http"

	reportsServer "github.com/pgquerynarrative/pgquerynarrative/api/gen/http/reports/server"
	"github.com/pgquerynarrative/pgquerynarrative/api/gen/reports"
	"github.com/pgquerynarrative/pgquerynarrative/app/llm"
	"github.com/pgquerynarrative/pgquerynarrative/app/queryrunner"
	"github.com/pgquerynarrative/pgquerynarrative/app/service"
	"github.com/pgquerynarrative/pgquerynarrative/test/testhelpers"
)

// mockLLM is used by reports E2E for List/Get only; Generate is not called.
type mockLLM struct{}

func (m *mockLLM) Generate(ctx context.Context, prompt string) (string, error) {
	return "", nil
}
func (m *mockLLM) Name() string { return "test" }

func TestReportsListAndGetE2E(t *testing.T) {
	ctx := context.Background()
	container := testhelpers.RunPostgresContainer(t, ctx)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	for attempt := 0; ; attempt++ {
		pool, pingErr := pgxpool.New(waitCtx, connStr)
		if pingErr == nil {
			pingErr = pool.Ping(waitCtx)
			pool.Close()
			if pingErr == nil {
				break
			}
		}
		if waitCtx.Err() != nil {
			t.Fatalf("postgres not ready after 15s: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	migrationsPath, err := filepath.Abs("../../app/db/migrations")
	if err != nil {
		t.Fatalf("failed to resolve migrations path: %v", err)
	}
	m, err := migrate.New("file://"+migrationsPath, connStr)
	if err != nil {
		t.Fatalf("failed to create migrator: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to run migrations: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	// Ensure app user can write to app.reports (migrations create schema and tables; we may need to grant)
	// Use the same connection string which is typically postgres superuser from testcontainers
	narrativeJSON := []byte(`{"headline":"E2E test report","takeaways":["One insight"],"drivers":[],"limitations":[],"recommendations":[]}`)
	metricsJSON := []byte(`{"aggregates":{},"data_quality":{},"time_series":{},"perf_suggestions":[]}`)

	var reportID string
	err = pool.QueryRow(ctx, `
		INSERT INTO app.reports (sql, narrative_md, narrative_json, metrics, llm_model, llm_provider)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, "SELECT 1", "E2E test report", narrativeJSON, metricsJSON, "test", "e2e").Scan(&reportID)
	if err != nil {
		t.Fatalf("failed to insert report: %v", err)
	}

	validator := queryrunner.NewValidator([]string{"demo"}, 10000)
	runner := queryrunner.NewRunner(pool, validator, 1000, 30*time.Second)
	var llmClient llm.Client = &mockLLM{}
	reportsService := service.NewReportsService(pool, pool, runner, llmClient, 0)
	endpoints := reports.NewEndpoints(reportsService)

	mux := goahttp.NewMuxer()
	dec := goahttp.RequestDecoder
	enc := goahttp.ResponseEncoder
	errHandler := func(ctx context.Context, w http.ResponseWriter, err error) {
		_ = goahttp.ErrorEncoder(enc, nil)(ctx, w, err)
	}
	httpSrv := reportsServer.New(endpoints, mux, dec, enc, errHandler, nil)
	reportsServer.Mount(mux, httpSrv)

	testServer := httptest.NewServer(mux)
	t.Cleanup(testServer.Close)

	// GET /api/v1/reports (list)
	listResp, err := http.Get(testServer.URL + "/api/v1/reports?limit=10&offset=0")
	if err != nil {
		t.Fatalf("list reports request failed: %v", err)
	}
	defer listResp.Body.Close()
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("list reports unexpected status: %d", listResp.StatusCode)
	}
	var listResult reports.ReportList
	if err := json.NewDecoder(listResp.Body).Decode(&listResult); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	if len(listResult.Items) == 0 {
		t.Fatal("expected at least one report in list")
	}
	if listResult.Items[0].ID != reportID {
		t.Errorf("list first item id = %q, want %q", listResult.Items[0].ID, reportID)
	}

	// GET /api/v1/reports/{id}
	getResp, err := http.Get(testServer.URL + "/api/v1/reports/" + reportID)
	if err != nil {
		t.Fatalf("get report request failed: %v", err)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("get report unexpected status: %d", getResp.StatusCode)
	}
	var report reports.Report
	if err := json.NewDecoder(getResp.Body).Decode(&report); err != nil {
		t.Fatalf("failed to decode get response: %v", err)
	}
	if report.ID != reportID {
		t.Errorf("get report id = %q, want %q", report.ID, reportID)
	}
	if report.SQL != "SELECT 1" {
		t.Errorf("get report sql = %q, want SELECT 1", report.SQL)
	}
	if report.Narrative == nil || report.Narrative.Headline != "E2E test report" {
		t.Errorf("get report narrative = %+v", report.Narrative)
	}

	// GET non-existent report -> 404
	notFoundResp, err := http.Get(testServer.URL + "/api/v1/reports/00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Fatalf("get missing report request failed: %v", err)
	}
	notFoundResp.Body.Close()
	if notFoundResp.StatusCode != http.StatusNotFound {
		t.Errorf("get missing report: want 404, got %d", notFoundResp.StatusCode)
	}
}
