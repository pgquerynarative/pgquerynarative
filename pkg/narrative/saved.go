package narrative

import (
	"context"

	"github.com/pgquerynarrative/pgquerynarrative/api/gen/queries"
)

// ListSavedQueries returns saved queries with optional limit and offset.
// Context cancellation is propagated.
func (c *Client) ListSavedQueries(ctx context.Context, limit, offset int) (*queries.SavedQueryList, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	payload := &queries.ListSavedPayload{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	return c.queriesService.ListSaved(ctx, payload)
}

// GetSavedQuery returns a saved query by ID. Returns queries.NotFoundError if not found.
// Context cancellation is propagated.
func (c *Client) GetSavedQuery(ctx context.Context, id string) (*queries.SavedQuery, error) {
	payload := &queries.GetSavedPayload{ID: id}
	return c.queriesService.GetSaved(ctx, payload)
}

// SaveQuery creates or updates a saved query. Name and SQL are required.
// Description and tags are optional (pass nil for description, nil or empty slice for tags).
// Context cancellation is propagated.
func (c *Client) SaveQuery(ctx context.Context, name, sql string, description *string, tags []string) (*queries.SavedQuery, error) {
	payload := &queries.SaveQueryPayload{
		Name: name,
		SQL:  sql,
	}
	if description != nil {
		payload.Description = description
	}
	if len(tags) > 0 {
		payload.Tags = tags
	}
	return c.queriesService.Save(ctx, payload)
}

// DeleteSavedQuery removes a saved query by ID. Returns queries.NotFoundError if not found.
// Context cancellation is propagated.
func (c *Client) DeleteSavedQuery(ctx context.Context, id string) error {
	payload := &queries.DeleteSavedPayload{ID: id}
	return c.queriesService.DeleteSaved(ctx, payload)
}
