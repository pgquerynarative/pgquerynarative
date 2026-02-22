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

	schemaServer "github.com/pgquerynarrative/pgquerynarrative/api/gen/http/schema/server"
	suggestionsServer "github.com/pgquerynarrative/pgquerynarrative/api/gen/http/suggestions/server"
	schema "github.com/pgquerynarrative/pgquerynarrative/api/gen/schema"
	suggestions "github.com/pgquerynarrative/pgquerynarrative/api/gen/suggestions"
	"github.com/pgquerynarrative/pgquerynarrative/app/catalog"
	"github.com/pgquerynarrative/pgquerynarrative/app/service"
	pkgsuggestions "github.com/pgquerynarrative/pgquerynarrative/app/suggestions"
	"github.com/pgquerynarrative/pgquerynarrative/test/testhelpers"
)

func TestSchemaAndSuggestionsE2E(t *testing.T) {
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

	_, err = pool.Exec(ctx, `
		INSERT INTO demo.sales (id, date, product_category, product_name, quantity, unit_price, total_amount, region, sales_rep)
		VALUES (gen_random_uuid(), CURRENT_DATE, 'Electronics', 'Alpha', 5, 10.00, 50.00, 'North', 'A. Lee')
	`)
	if err != nil {
		t.Fatalf("failed to seed data: %v", err)
	}

	loader := catalog.NewLoader(pool, []string{"demo"})
	schemaService := service.NewSchemaService(loader)
	suggester := pkgsuggestions.NewSuggester(pool)

	schemaEndpoints := schema.NewEndpoints(schemaService)
	suggestionsEndpoints := suggestions.NewEndpoints(suggester)

	mux := goahttp.NewMuxer()
	dec := goahttp.RequestDecoder
	enc := goahttp.ResponseEncoder
	errHandler := func(ctx context.Context, w http.ResponseWriter, err error) {
		_ = goahttp.ErrorEncoder(enc, nil)(ctx, w, err)
	}

	schemaHTTP := schemaServer.New(schemaEndpoints, mux, dec, enc, errHandler, nil)
	schemaServer.Mount(mux, schemaHTTP)
	suggestionsHTTP := suggestionsServer.New(suggestionsEndpoints, mux, dec, enc, errHandler, nil)
	suggestionsServer.Mount(mux, suggestionsHTTP)

	testServer := httptest.NewServer(mux)
	t.Cleanup(testServer.Close)

	// GET /api/v1/schema
	schemaResp, err := http.Get(testServer.URL + "/api/v1/schema")
	if err != nil {
		t.Fatalf("schema request failed: %v", err)
	}
	defer schemaResp.Body.Close()
	if schemaResp.StatusCode != http.StatusOK {
		t.Fatalf("schema unexpected status: %d", schemaResp.StatusCode)
	}
	var schemaResult schema.SchemaResult
	if err := json.NewDecoder(schemaResp.Body).Decode(&schemaResult); err != nil {
		t.Fatalf("failed to decode schema response: %v", err)
	}
	if len(schemaResult.Schemas) == 0 {
		t.Fatal("expected at least one schema (demo)")
	}
	var demoFound bool
	for _, s := range schemaResult.Schemas {
		if s.Name == "demo" {
			demoFound = true
			if len(s.Tables) == 0 {
				t.Fatal("expected demo to have tables")
			}
			break
		}
	}
	if !demoFound {
		t.Errorf("expected demo schema, got: %v", schemaResult.Schemas)
	}

	// GET /api/v1/suggestions/queries
	suggResp, err := http.Get(testServer.URL + "/api/v1/suggestions/queries?limit=5")
	if err != nil {
		t.Fatalf("suggestions request failed: %v", err)
	}
	defer suggResp.Body.Close()
	if suggResp.StatusCode != http.StatusOK {
		t.Fatalf("suggestions unexpected status: %d", suggResp.StatusCode)
	}
	var suggResult suggestions.SuggestedQueriesResult
	if err := json.NewDecoder(suggResp.Body).Decode(&suggResult); err != nil {
		t.Fatalf("failed to decode suggestions response: %v", err)
	}
	if len(suggResult.Suggestions) == 0 {
		t.Fatal("expected at least one suggestion (curated)")
	}
}
