package narrative

// Public types for the library are aligned with the Goa API. RunQuery and
// GenerateReport return the same payload/result types as the HTTP API so that
// library and server stay in sync. Callers should use:
//
//   - queries.RunQueryPayload / queries.RunQueryResult for RunQuery
//   - reports.GenerateReportPayload / reports.Report for GenerateReport
//
// These are defined in api/gen/queries and api/gen/reports. This package
// does not duplicate them; import those packages when you need the concrete
// types for responses.

// RunQueryOptions holds optional settings for query execution. Zero value uses defaults.
type RunQueryOptions struct {
	// Limit is the maximum number of rows to return. If <= 0, a default (e.g. 1000) is used.
	Limit int
}

// DefaultRunQueryLimit is the limit used when RunQuery is called with limit <= 0.
const DefaultRunQueryLimit = 1000
