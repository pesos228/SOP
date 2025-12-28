package poolgrp

import (
	"fmt"
	"hosting-kit/page"
	"hosting-resources-service/cmd/server/rest/gen"
	"hosting-resources-service/cmd/server/rest/pagination"
	"hosting-resources-service/internal/pool"
)

func toPool(p pool.Pool, prefix string) gen.Pool {
	links := gen.Links{
		"self": gen.Link{Href: fmt.Sprintf("%s/plans/%s", prefix, p.ID)},
	}

	resources := gen.Resource{
		CpuCores: p.Resources.CPUCores,
		DiskGb:   p.Resources.DiskGB,
		IpCount:  p.Resources.IPCount,
		RamMb:    p.Resources.RAMMB,
	}

	return gen.Pool{
		UnderscoreLinks: links,
		Id:              p.ID,
		Name:            p.Name,
		Resources:       resources,
	}
}

func toPoolCollectionResponse(pools []pool.Pool, pg page.Page, total int, prefix string) gen.PoolCollectionResponse {
	items := make([]gen.Pool, len(pools))
	for i, p := range pools {
		items[i] = toPool(p, prefix)
	}

	return gen.PoolCollectionResponse{
		UnderscoreEmbedded: struct {
			Pools []gen.Pool `json:"pools"`
		}{
			Pools: items,
		},
		Page:            pagination.ToMetaData(pg, total),
		UnderscoreLinks: pagination.ToLinks(fmt.Sprintf("%s/pools", prefix), pg, total),
	}
}
