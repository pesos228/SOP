-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pools (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    cpu_cores INT NOT NULL,
    ram_mb INT NOT NULL,
    disk_gb INT NOT NULL,
    ip_count INT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_pools_updated_at ON pools(updated_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pools;
-- +goose StatementEnd
