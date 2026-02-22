package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pgquerynarrative/pgquerynarrative/pkg/narrative"
)

// MountGin mounts PgQueryNarrative HTTP handlers on the Gin engine under prefix.
// Routes: POST prefix/query/run, POST prefix/report/generate, GET prefix/schema,
// GET prefix/suggestions/queries. Prefix is normalized (no trailing slash).
// Client must not be nil.
func MountGin(r *gin.Engine, client *narrative.Client, prefix string) {
	prefix = strings.TrimSuffix(prefix, "/")
	g := r.Group(prefix)
	g.POST("/query/run", wrapGin(client, HandleRunQuery))
	g.POST("/report/generate", wrapGin(client, HandleGenerateReport))
	g.GET("/schema", wrapGin(client, HandleGetSchema))
	g.GET("/suggestions/queries", wrapGin(client, HandleSuggestionsQueries))
}

func wrapGin(client *narrative.Client, fn func(*narrative.Client, http.ResponseWriter, *http.Request)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(client, c.Writer, c.Request)
	}
}
