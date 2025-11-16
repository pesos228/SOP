package app

import (
	"context"
	"hosting-kit/messaging"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"hosting-service/internal/config"
	"hosting-service/internal/repository/psql"
	"hosting-service/internal/service"

	"gorm.io/gorm"
)

type App struct {
	config *config.Config
}

func New(cfg *config.Config) *App {
	return &App{config: cfg}
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

	var wg sync.WaitGroup

	rabbit, err := a.initRabbitManager(&wg)
	if err != nil {
		log.Fatalf("RabbitMQ initialization failed: %v", err)
	}

	services := buildDependencies(db, rabbit)

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a.runConsumers(rabbit, services.serverService)

	httpServer := a.runHTTPServer(&wg, services)

	log.Printf("Сервер GraphQL запущен. Playground: http://localhost%s/graphi\n", a.config.HTTP_Port)
	log.Printf("Сервер REST (Swagger) запущен. UI: http://localhost%s/swagger/\n", a.config.HTTP_Port)

	<-shutdownCtx.Done()
	log.Println("Shutting down application gracefully...")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Closing HTTP server...")

	if err := httpServer.Shutdown(timeoutCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Closing RabbitMQ consumer...")

	rabbit.Stop(timeoutCtx)

	wg.Wait()
}

func buildDependencies(db *gorm.DB, rabbit *messaging.MessageManager) services {
	planRepository := psql.NewPlanRepository(db)
	serverRepository := psql.NewServerRepository(db)

	planService := service.NewPlanService(planRepository)
	serverService := service.NewServerService(serverRepository, planRepository, rabbit)

	return services{
		planService:   planService,
		serverService: serverService,
	}
}
