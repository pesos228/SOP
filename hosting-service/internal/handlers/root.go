package handlers

import (
	"encoding/json"
	"hosting-service/internal/assemblers"
	"hosting-service/internal/dto"
	"hosting-service/internal/service"
	"net/http"
)

type ApiHandler struct {
	*PlanHandler
	*ServersHandler
}

func NewApiHandler(planService service.PlanService, serverService service.ServerService) *ApiHandler {
	return &ApiHandler{
		PlanHandler:    NewPlansHandler(planService),
		ServersHandler: NewServersHandler(serverService),
	}
}

func RootEntryPoint(w http.ResponseWriter, r *http.Request) {
	var routes []dto.EntryPointLink

	routes = append(routes,
		dto.EntryPointLink{
			Rel:       "plans",
			Href:      "/api/plans{?page,pageSize}",
			Templated: true,
		},
		dto.EntryPointLink{
			Rel:       "servers",
			Href:      "/api/servers{?page,pageSize}",
			Templated: true,
		},
		dto.EntryPointLink{
			Rel:  "documentation",
			Href: "/swagger/index.html",
		},
	)

	w.Header().Set("Content-Type", "application/hal+json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(assemblers.ToRouteCollectionResponse(routes))
}
