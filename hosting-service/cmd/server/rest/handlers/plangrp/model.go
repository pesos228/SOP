package plangrp

import (
	"fmt"
	"hosting-kit/page"
	"hosting-service/cmd/server/rest/gen"
	"hosting-service/cmd/server/rest/pagination"
	"hosting-service/internal/plan"
)

func toPlan(p plan.Plan, prefix string) gen.ServerPlan {
	links := gen.Links{
		"self": gen.Link{Href: fmt.Sprintf("%s/plans/%s", prefix, p.ID)},
	}

	return gen.ServerPlan{
		Id:              p.ID,
		Name:            p.Name,
		CpuCores:        p.CPUCores,
		RamMb:           p.RAMMB,
		DiskGb:          p.DiskGB,
		UnderscoreLinks: links,
	}
}

func toPlanCollectionResponse(plans []plan.Plan, pg page.Page, total int, prefix string) gen.PlanCollectionResponse {
	items := make([]gen.ServerPlan, len(plans))
	for i, p := range plans {
		items[i] = toPlan(p, prefix)
	}

	return gen.PlanCollectionResponse{
		UnderscoreEmbedded: struct {
			Plans []gen.ServerPlan `json:"plans"`
		}{
			Plans: items,
		},
		Page:            pagination.ToMetaData(pg, total),
		UnderscoreLinks: pagination.ToLinks(fmt.Sprintf("%s/plans", prefix), pg, total),
	}
}
