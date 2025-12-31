-- +goose Up
-- +goose StatementBegin
ALTER TABLE plans ADD COLUMN ip_count INT NOT NULL DEFAULT 1;
ALTER TABLE servers ADD COLUMN pool_id UUID NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE servers DROP COLUMN pool_id;
ALTER TABLE plans DROP COLUMN ip_count;
-- +goose StatementEnd
