package app

import (
	events "hosting-events-contract"
	"hosting-kit/messaging"
	"log"
	"sync"

	"hosting-service/internal/listeners"
	"hosting-service/internal/service"
)

func (a *App) initRabbitManager(wg *sync.WaitGroup) (*messaging.MessageManager, error) {
	exchanges := []messaging.ExchangeConfig{
		{
			Name: events.CommandsExchange,
			Type: messaging.ExchangeDirect,
		},
		{
			Name: events.EventsExchange,
			Type: messaging.ExchangeTopic,
		},
	}
	return messaging.NewMessageManager(a.config.AMQP_URL, exchanges, wg, a.config.AMQP_HandlerTimeout)
}

func (a *App) runConsumers(rabbit *messaging.MessageManager, serverService service.ServerService) {
	resultsListener := listeners.NewProvisioningResultListener(serverService)

	err := rabbit.Subscribe(
		a.config.ResultsQueue,
		events.ProvisionSucceededKey,
		events.EventsExchange,
		messaging.ExchangeTopic,
		resultsListener.HandleSuccess,
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	err = rabbit.Subscribe(
		a.config.ResultsQueue,
		events.ProvisionFailedKey,
		events.EventsExchange,
		messaging.ExchangeTopic,
		resultsListener.HandleFailure,
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
}
