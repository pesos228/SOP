package pooldb

import (
	"hosting-resources-service/internal/pool"
	"time"

	"github.com/google/uuid"
)

type poolDB struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CPUCores  int       `db:"cpu_cores"`
	RAMMB     int       `db:"ram_mb"`
	DiskGB    int       `db:"disk_gb"`
	IPCount   int       `db:"ip_count"`
	UpdatedAt time.Time `db:"updated_at"`
}

func toDBPool(p pool.Pool) poolDB {
	return poolDB{
		ID:        p.ID,
		Name:      p.Name,
		CPUCores:  p.Resources.CPUCores,
		RAMMB:     p.Resources.RAMMB,
		DiskGB:    p.Resources.DiskGB,
		IPCount:   p.Resources.IPCount,
		UpdatedAt: time.Now().UTC(),
	}
}

func toBusPool(db poolDB) pool.Pool {
	res := pool.Resource{
		CPUCores: db.CPUCores,
		RAMMB:    db.RAMMB,
		DiskGB:   db.DiskGB,
		IPCount:  db.IPCount,
	}

	return pool.Pool{
		ID:        db.ID,
		Name:      db.Name,
		Resources: res,
	}
}

func toBusPools(dbs []poolDB) []pool.Pool {
	pools := make([]pool.Pool, len(dbs))
	for i, db := range dbs {
		pools[i] = toBusPool(db)
	}
	return pools
}
