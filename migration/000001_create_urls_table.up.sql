CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    target TEXT NOT NULL,
    code VARCHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_urls_id_desc ON urls (id DESC);