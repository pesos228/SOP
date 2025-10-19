package assemblers

import (
	"hosting-contracts/api"
	"hosting-service/internal/dto"
)

func toRoute(route string, templated bool) api.Link {
	return api.Link{
		Href:      to(route),
		Templated: to(templated),
	}
}

func ToRouteCollectionResponse(routes []dto.EntryPointLink) api.Links {
	links := make(api.Links)

	for _, route := range routes {
		links[route.Rel] = toRoute(route.Href, route.Templated)
	}

	return links
}
