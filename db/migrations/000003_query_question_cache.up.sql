CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE query_question_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    sql_query TEXT NOT NULL,
    embedding vector(1536) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_query_question_cache_user_id ON query_question_cache (user_id);

CREATE INDEX idx_query_question_cache_embedding
    ON query_question_cache
    USING hnsw (embedding vector_cosine_ops);
