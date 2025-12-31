package pooldb

import (
	"context"
	"errors"
	"fmt"
	"hosting-kit/page"
	"hosting-resources-service/internal/pool"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db}
}

func (s *Store) AppendResource(ctx context.Context, r pool.Resource, poolID uuid.UUID) (pool.Pool, error) {
	const q = `
	UPDATE pools
	SET
		cpu_cores  = cpu_cores + @cpu,
		ram_mb     = ram_mb    + @ram,
		disk_gb    = disk_gb   + @disk,
		ip_count   = ip_count  + @ip,
		updated_at = NOW()
	WHERE
		id = @id
	RETURNING id, name, cpu_cores, ram_mb, disk_gb, ip_count, updated_at`

	args := pgx.NamedArgs{
		"id":   poolID,
		"cpu":  r.CPUCores,
		"ram":  r.RAMMB,
		"disk": r.DiskGB,
		"ip":   r.IPCount,
	}

	var dbPool poolDB
	err := s.db.QueryRow(ctx, q, args).Scan(
		&dbPool.ID,
		&dbPool.Name,
		&dbPool.CPUCores,
		&dbPool.RAMMB,
		&dbPool.DiskGB,
		&dbPool.IPCount,
		&dbPool.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pool.Pool{}, pool.ErrPoolNotFound
		}
		return pool.Pool{}, fmt.Errorf("db: append query: %w", err)
	}

	return toBusPool(dbPool), nil
}

func (s *Store) SubtractResource(ctx context.Context, r pool.Resource) (uuid.UUID, error) {
	const q = `
	UPDATE pools
	SET
		cpu_cores  = cpu_cores - @cpu,
		ram_mb     = ram_mb    - @ram,
		disk_gb    = disk_gb   - @disk,
		ip_count   = ip_count  - @ip,
		updated_at = NOW()
	WHERE id = (
		SELECT id
		FROM pools
		WHERE 
			cpu_cores >= @cpu AND
			ram_mb    >= @ram AND
			disk_gb   >= @disk AND
			ip_count  >= @ip
		ORDER BY updated_at ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	)
	RETURNING id`

	args := pgx.NamedArgs{
		"cpu":  r.CPUCores,
		"ram":  r.RAMMB,
		"disk": r.DiskGB,
		"ip":   r.IPCount,
	}

	var poolID uuid.UUID
	err := s.db.QueryRow(ctx, q, args).Scan(&poolID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, pool.ErrNotEnoughResources
		}
		return uuid.Nil, fmt.Errorf("db: subtract exec: %w", err)
	}

	return poolID, nil
}

func (s *Store) CreatePool(ctx context.Context, p pool.Pool) error {
	const q = `
	INSERT INTO pools
		(id, name, cpu_cores, ram_mb, disk_gb, ip_count, updated_at)
	VALUES
		(@id, @name, @cpu_cores, @ram_mb, @disk_gb, @ip_count, @updated_at)
	`

	dbPool := toDBPool(p)

	args := pgx.NamedArgs{
		"id":         dbPool.ID,
		"name":       dbPool.Name,
		"cpu_cores":  dbPool.CPUCores,
		"disk_gb":    dbPool.DiskGB,
		"ram_mb":     dbPool.RAMMB,
		"ip_count":   dbPool.IPCount,
		"updated_at": dbPool.UpdatedAt,
	}

	_, err := s.db.Exec(ctx, q, args)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}

	return nil
}

func (s *Store) FindAll(ctx context.Context, pg page.Page) ([]pool.Pool, int, error) {
	const qCount = `SELECT count(*) FROM pools`

	var total int
	if err := s.db.QueryRow(ctx, qCount).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("db: %w", err)
	}

	if total == 0 {
		return []pool.Pool{}, 0, nil
	}

	const qSelect = `
	SELECT 
		id, name, cpu_cores, ram_mb, disk_gb, ip_count, updated_at
	FROM 
		pools
	ORDER BY 
		id ASC
	LIMIT @limit OFFSET @offset`

	args := pgx.NamedArgs{
		"limit":  pg.Size(),
		"offset": pg.Offset(),
	}

	rows, err := s.db.Query(ctx, qSelect, args)
	if err != nil {
		return nil, 0, fmt.Errorf("db: %w", err)
	}

	dbPools, err := pgx.CollectRows(rows, pgx.RowToStructByName[poolDB])
	if err != nil {
		return nil, 0, fmt.Errorf("db: %w", err)
	}

	return toBusPools(dbPools), total, nil
}
