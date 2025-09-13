CREATE TABLE urls (
    id              SERIAL PRIMARY KEY,
    code            VARCHAR(32) NOT NULL UNIQUE,
    original        TEXT NOT NULL,
    click_count     INTEGER NOT NULL DEFAULT 0,
    expiration_date TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '1 YEAR'),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_urls_code ON urls(code);
CREATE INDEX idx_urls_expiration_date ON urls(expiration_date);