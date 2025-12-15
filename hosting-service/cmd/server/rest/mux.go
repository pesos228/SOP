package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"hosting-contracts/api"
	"hosting-service/internal/plan"
	"hosting-service/internal/server"
)

type Config struct {
	PlanBus   *plan.Business
	ServerBus *server.Business
}

func RegisterRoutes(router *chi.Mux, cfg Config) {
	apiImpl := New(cfg.PlanBus, cfg.ServerBus)
	strictHandler := api.NewStrictHandler(apiImpl, nil)

	router.Route("/api", func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/api/swagger/doc.yaml"),
		))
		r.Get("/swagger/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-yaml")
			w.Write(api.OpenApiSpec)
		})

		api.HandlerFromMux(strictHandler, r)
	})
}
