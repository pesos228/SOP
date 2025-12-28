package pool

import "github.com/google/uuid"

type Pool struct {
	ID        uuid.UUID
	Name      string
	Resources Resource
}

type Resource struct {
	CPUCores int
	RAMMB    int
	DiskGB   int
	IPCount  int
}

type NewPool struct {
	Name     string
	CPUCores int
	RAMMB    int
	DiskGB   int
	IPCount  int
}
