package servergrp

import (
	"fmt"
	"hosting-kit/page"
	"hosting-service/cmd/server/rest/gen"
	"hosting-service/cmd/server/rest/pagination"
	"hosting-service/internal/server"
)

func toServer(s server.Server, prefix string) gen.Server {
	links := make(gen.Links)

	selfLink := fmt.Sprintf("%s/servers/%s", prefix, s.ID)
	actionsLink := fmt.Sprintf("%s/servers/%s/actions", prefix, s.ID)

	links["self"] = gen.Link{Href: selfLink}

	switch s.Status {
	case server.StatusRunning:
		links["stop"] = gen.Link{Href: actionsLink}
	case server.StatusStopped:
		links["start"] = gen.Link{Href: actionsLink}
		links["delete"] = gen.Link{Href: actionsLink}
	}

	return gen.Server{
		Id:              s.ID,
		Name:            s.Name,
		PlanId:          s.PlanID,
		IPv4Address:     s.IPv4Address,
		Status:          gen.ServerStatus(s.Status),
		CreatedAt:       s.CreatedAt,
		UnderscoreLinks: links,
	}
}

func toServerCollectionResponse(servers []server.Server, pg page.Page, total int, prefix string) gen.ServerCollectionResponse {
	items := make([]gen.Server, len(servers))
	for i, s := range servers {
		items[i] = toServer(s, prefix)
	}

	return gen.ServerCollectionResponse{
		UnderscoreEmbedded: struct {
			Servers []gen.Server `json:"servers"`
		}{
			Servers: items,
		},
		Page:            pagination.ToMetaData(pg, total),
		UnderscoreLinks: pagination.ToLinks(fmt.Sprintf("%s/servers", prefix), pg, total),
	}
}
