package embedding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SimilarQuery holds a saved query and its similarity score (0–1, higher is more similar).
type SimilarQuery struct {
	SavedQueryID string
	Name         string
	SQL          string
	Description  string
	Score        float64
}

type scoredQuery struct {
	sim   SimilarQuery
	score float64
}

// Store persists and retrieves query embeddings for similar-query search and RAG.
type Store struct {
	appPool *pgxpool.Pool
}

// NewStore creates a store that uses the app pool (writes to app.query_embeddings).
func NewStore(appPool *pgxpool.Pool) *Store {
	return &Store{appPool: appPool}
}

// Upsert saves or replaces the embedding for a saved query. Embedding is stored as JSONB.
func (s *Store) Upsert(ctx context.Context, savedQueryID string, embedding []float32, model string) error {
	raw, err := json.Marshal(embedding)
	if err != nil {
		return fmt.Errorf("marshal embedding: %w", err)
	}
	_, err = s.appPool.Exec(ctx, `
		INSERT INTO app.query_embeddings (saved_query_id, embedding, model, updated_at)
		VALUES ($1::uuid, $2::jsonb, $3, NOW())
		ON CONFLICT (saved_query_id) DO UPDATE SET embedding = EXCLUDED.embedding, model = EXCLUDED.model, updated_at = NOW()
	`, savedQueryID, raw, model)
	if err != nil {
		return fmt.Errorf("upsert embedding: %w", err)
	}
	return nil
}

// FindSimilar returns saved queries most similar to the given embedding (cosine similarity).
// Loads all stored embeddings and ranks in memory; limit is the max number to return.
func (s *Store) FindSimilar(ctx context.Context, queryEmbedding []float32, limit int) ([]SimilarQuery, error) {
	if limit <= 0 {
		limit = 5
	}
	rows, err := s.appPool.Query(ctx, `
		SELECT qe.saved_query_id::text, qe.embedding, sq.name, sq.sql, COALESCE(sq.description, '')
		FROM app.query_embeddings qe
		JOIN app.saved_queries sq ON sq.id = qe.saved_query_id
	`)
	if err != nil {
		return nil, fmt.Errorf("query embeddings: %w", err)
	}
	defer rows.Close()

	type row struct {
		id            string
		embeddingJSON []byte
		name          string
		sql           string
		description   string
	}
	var candidates []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.id, &r.embeddingJSON, &r.name, &r.sql, &r.description); err != nil {
			return nil, err
		}
		candidates = append(candidates, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cosine similarity: score = dot(a,b) when vectors are L2-normalized (Ollama returns normalized).
	queryNorm := norm(queryEmbedding)
	if queryNorm == 0 {
		return nil, nil
	}
	var scoredList []scoredQuery
	for _, c := range candidates {
		var vec []float32
		if err := json.Unmarshal(c.embeddingJSON, &vec); err != nil {
			continue
		}
		score := cosineSimilarity(queryEmbedding, vec)
		scoredList = append(scoredList, scoredQuery{
			sim: SimilarQuery{
				SavedQueryID: c.id,
				Name:         c.name,
				SQL:          c.sql,
				Description:  c.description,
				Score:        score,
			},
			score: score,
		})
	}
	// Sort by score descending and take top limit
	sortByScoreDesc(scoredList)
	out := make([]SimilarQuery, 0, limit)
	for i := 0; i < len(scoredList) && i < limit; i++ {
		out = append(out, scoredList[i].sim)
	}
	return out, nil
}

func norm(v []float32) float64 {
	var sum float64
	for _, x := range v {
		sum += float64(x) * float64(x)
	}
	if sum <= 0 {
		return 0
	}
	return sqrt(sum)
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	// Newton step
	y := x
	for i := 0; i < 10; i++ {
		next := (y + x/y) / 2
		if next == y {
			return y
		}
		y = next
	}
	return y
}

func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
	}
	na := norm(a)
	nb := norm(b)
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (na * nb)
}

func sortByScoreDesc(list []scoredQuery) {
	// Simple insertion sort for small N
	for i := 1; i < len(list); i++ {
		for j := i; j > 0 && list[j].score > list[j-1].score; j-- {
			list[j], list[j-1] = list[j-1], list[j]
		}
	}
}
