package narrative

import (
	"context"

	"github.com/pgquerynarrative/pgquerynarrative/api/gen/queries"
	"github.com/pgquerynarrative/pgquerynarrative/api/gen/reports"
)

// GenerateReport runs the given SQL query and generates a narrative report
// (headline, takeaways, drivers, etc.) using the configured LLM. Context
// cancellation is propagated. The result is the same type as the API generate endpoint.
func (c *Client) GenerateReport(ctx context.Context, sql string) (*reports.Report, error) {
	payload := &reports.GenerateReportPayload{SQL: sql}
	return c.reportsService.Generate(ctx, payload)
}

// GenerateReportFromSaved generates a report from a saved query by ID. It fetches
// the saved query's SQL and then generates the report, associating the report
// with that saved query. Returns an error if the saved query is not found.
// Context cancellation is propagated.
func (c *Client) GenerateReportFromSaved(ctx context.Context, savedQueryID string) (*reports.Report, error) {
	saved, err := c.queriesService.GetSaved(ctx, &queries.GetSavedPayload{ID: savedQueryID})
	if err != nil {
		return nil, err
	}
	payload := &reports.GenerateReportPayload{
		SQL:          saved.SQL,
		SavedQueryID: &savedQueryID,
	}
	return c.reportsService.Generate(ctx, payload)
}
