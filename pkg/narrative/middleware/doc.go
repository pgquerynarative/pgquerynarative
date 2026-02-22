// Package middleware provides framework-specific helpers to mount PgQueryNarrative
// HTTP endpoints into an existing Go server (Chi, Gin, Echo). Use this for
// embedded integration: create a narrative.Client, then call MountChi, MountGin,
// or MountEcho to register query/run, report/generate, schema, and suggestions
// routes under a path prefix.
package middleware
