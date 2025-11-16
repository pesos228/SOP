package app

import (
	"fmt"
	"hosting-contracts/api"
	"log"
	"net/http"
	"sync"

	"hosting-service/internal/graph"
	"hosting-service/internal/handlers"
	"hosting-service/internal/middleware"
	"hosting-service/internal/service"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (a *App) runHTTPServer(wg *sync.WaitGroup, services services) *http.Server {
	router := a.setupRouter(services.planService, services.serverService)
	srv := &http.Server{Addr: a.config.HTTP_Port, Handler: router}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Starting HTTP server on %s", a.config.HTTP_Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server listen error: %s\n", err)
		}
		log.Println("HTTP server stopped.")
	}()

	return srv
}

func (a *App) setupRouter(planService service.PlanService, serverService service.ServerService) *chi.Mux {
	graphqlResolver := &graph.Resolver{
		PlanService:   planService,
		ServerService: serverService,
	}
	executableSchema := graph.NewExecutableSchema(graph.Config{Resolvers: graphqlResolver})
	graphqlHandler := handler.NewDefaultServer(executableSchema)
	playgroundHandler := playground.Handler("GraphQL Playground", "/graphql")

	apiHandler := handlers.NewApiHandler(planService, serverService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.PerformanceLogger)
	router.Use(chiMiddleware.Recoverer)

	router.Handle("/graphi", playgroundHandler)
	router.Handle("/graphql", graphqlHandler)

	router.Get("/swagger/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		if _, err := w.Write(api.OpenApiSpec); err != nil {
			log.Printf("Error writing OpenAPI spec: %v", err)
		}
	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost%s/swagger/doc.yaml", a.config.HTTP_Port)),
	))

	apiRouter := chi.NewRouter()
	strictHandler := api.NewStrictHandler(apiHandler, nil)
	api.HandlerFromMux(strictHandler, apiRouter)
	apiRouter.Get("/", handlers.RootEntryPoint)
	router.Mount("/api", apiRouter)

	return router
}
