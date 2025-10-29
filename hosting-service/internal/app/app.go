package app

import (
	"context"
	"fmt"
	"hosting-service/internal/config"
	"hosting-service/internal/domain"
	"hosting-service/internal/graph"
	"hosting-service/internal/handlers"
	"hosting-service/internal/listeners"
	"hosting-service/internal/messaging"
	"hosting-service/internal/middleware"
	"hosting-service/internal/repository/psql"
	"hosting-service/internal/service"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"hosting-contracts/api"

	events "hosting-events-contract"

	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	config *config.Config
}

func New(cfg *config.Config) *App {
	return &App{config: cfg}
}

func (a *App) Run() {
	db, err := a.initDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	conn, err := a.initRabbitMQ()
	if err != nil {
		log.Fatalf("RabbitMQ initialization failed: %v", err)
	}

	planRepo := psql.NewPlanRepository(db)
	serverRepo := psql.NewServerRepository(db)

	cmdPublisher, err := messaging.NewEventPublisher(conn, events.CommandsExchange)
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}

	planService := service.NewPlanService(planRepo)
	serverService := service.NewServerService(serverRepo, planRepo, cmdPublisher)

	var wg sync.WaitGroup
	resultsListener := listeners.NewProvisioningResultListener(serverService)
	resultsConsumer, err := messaging.NewManagedConsumer(
		conn,
		a.config.ResultsQueue,
		events.ProvisionResultKeyPattern,
		events.EventsExchange,
		resultsListener.Handle,
	)
	if err != nil {
		log.Fatalf("Failed to create results consumer: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting results listener...")
		if err := resultsConsumer.Run(); err != nil {
			log.Printf("Results listener stopped: %v", err)
		}
	}()

	router := a.setupRouter(planService, serverService)
	srv := &http.Server{Addr: a.config.HTTP_Port, Handler: router}

	go func() {
		log.Printf("Starting HTTP server on %s", a.config.HTTP_Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server listen error: %s\n", err)
		}
	}()

	log.Printf("Сервер GraphQL запущен. Playground: http://localhost%s/graphi\n", a.config.HTTP_Port)
	log.Printf("Сервер REST (Swagger) запущен. UI: http://localhost%s/swagger/\n", a.config.HTTP_Port)

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-shutdownCtx.Done()
	log.Println("Shutting down application gracefully...")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Closing RabbitMQ consumer...")

	consumerShutdownCtx, consumerCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer consumerCancel()
	resultsConsumer.CloseWithContext(consumerShutdownCtx)

	log.Println("Waiting for background goroutines to finish...")
	wg.Wait()

	log.Println("Closing publisher and main connection...")
	cmdPublisher.Close()
	conn.Close()

	log.Println("Application stopped successfully.")
}

func (a *App) initDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(a.config.DB_DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connection to DB failed: %w", err)
	}

	log.Println("Running DB migrations...")
	db.AutoMigrate(&domain.Plan{}, &domain.Server{})

	return db, nil
}

func (a *App) initRabbitMQ() (*rabbitmq.Conn, error) {
	conn, err := rabbitmq.NewConn(a.config.AMQP_URL, rabbitmq.WithConnectionOptionsLogging)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return conn, nil
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
		w.Write(api.OpenApiSpec)
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
