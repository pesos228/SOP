-- +goose Up
-- +goose StatementBegin
INSERT INTO pools (id, name, cpu_cores, ram_mb, disk_gb, ip_count, updated_at)
VALUES
(
    '11111111-1111-1111-1111-111111111111',
    'General Purpose Pool',
    100,
    256000,
    10000,
    50,
    NOW()
),
(
    '22222222-2222-2222-2222-222222222222',
    'High Performance Pool',
    500,
    1024000,
    50000,
    200,
    NOW()
)
ON CONFLICT (id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM pools WHERE id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222'
);
-- +goose StatementEnd
