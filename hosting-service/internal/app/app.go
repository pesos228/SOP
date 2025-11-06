package app

import (
	"context"
	"hosting-service/internal/config"
	"hosting-service/internal/messaging"
	"hosting-service/internal/repository"
	"hosting-service/internal/repository/psql"
	"hosting-service/internal/service"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"

	events "hosting-events-contract"
)

type App struct {
	config *config.Config
}

func New(cfg *config.Config) *App {
	return &App{config: cfg}
}

type repositories struct {
	planRepository   repository.PlanRepository
	serverRepository repository.ServerRepository
}

type services struct {
	planService   service.PlanService
	serverService service.ServerService
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

	publisher, _, services := buildDependencies(db, conn)
	defer publisher.Close()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	a.runMessageConsumers(shutdownCtx, &wg, conn, services.serverService)

	httpServer := a.runHTTPServer(&wg, services)

	log.Printf("Сервер GraphQL запущен. Playground: http://localhost%s/graphi\n", a.config.HTTP_Port)
	log.Printf("Сервер REST (Swagger) запущен. UI: http://localhost%s/swagger/\n", a.config.HTTP_Port)

	<-shutdownCtx.Done()
	log.Println("Shutting down application gracefully...")

	shutdownHTTPServer(httpServer)

	log.Println("Closing RabbitMQ consumer...")

	wg.Wait()
}

func buildDependencies(db *gorm.DB, conn *rabbitmq.Conn) (*messaging.EventPublisher, repositories, services) {
	planRepository := psql.NewPlanRepository(db)
	serverRepository := psql.NewServerRepository(db)

	cmdPublisher, err := messaging.NewEventPublisher(conn, events.CommandsExchange)
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}

	planService := service.NewPlanService(planRepository)
	serverService := service.NewServerService(serverRepository, planRepository, cmdPublisher)

	return cmdPublisher, repositories{
			planRepository:   planRepository,
			serverRepository: serverRepository,
		},
		services{
			planService:   planService,
			serverService: serverService,
		}
}
