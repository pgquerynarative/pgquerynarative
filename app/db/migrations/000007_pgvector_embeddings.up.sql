-- Optional: enable pgvector for in-database semantic search (vectorization inside Postgres).
-- Requires: PostgreSQL with pgvector installed and the vector extension created by a superuser.
-- Run once as superuser: ./tools/db/ensure-pgvector-extension.sh  (or: psql -d pgquerynarrative -c 'CREATE EXTENSION IF NOT EXISTS vector;')
-- If the extension is missing, this migration will fail; the app then uses in-memory similarity instead.

ALTER TABLE app.query_embeddings ADD COLUMN IF NOT EXISTS embedding_vector vector(768);

CREATE INDEX IF NOT EXISTS idx_query_embeddings_vector_cosine
ON app.query_embeddings USING hnsw (embedding_vector vector_cosine_ops)
WITH (m = 16, ef_construction = 64);
