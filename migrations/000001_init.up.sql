CREATE TABLE IF NOT EXISTS news (
    id BIGSERIAL PRIMARY KEY,
    url_hash TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL,
    title TEXT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    source TEXT NOT NULL,
    topics TEXT[] NOT NULL DEFAULT '{}',
    published_at TIMESTAMPTZ,
    language TEXT NOT NULL DEFAULT 'en',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_news_published_at ON news (published_at DESC);
CREATE INDEX IF NOT EXISTS idx_news_source ON news (source);
CREATE INDEX IF NOT EXISTS idx_news_topics ON news USING GIN (topics);
