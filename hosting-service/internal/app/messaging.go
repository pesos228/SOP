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
		{
			Name: events.DLXExchange,
			Type: messaging.ExchangeDirect,
		},
	}
	return messaging.NewMessageManager(a.config.AMQP_URL, exchanges, wg, a.config.AMQP_HandlerTimeout)
}

func (a *App) runConsumers(rabbit *messaging.MessageManager, serverService service.ServerService) {
	resultsListener := listeners.NewProvisioningResultListener(serverService)
	deadLetterListener := listeners.NewDeadLetterListener()

	err := rabbit.Subscribe(
		a.config.ResultsQueue,
		events.ProvisionResultKeyPattern,
		events.EventsExchange,
		resultsListener.Handle,
		&messaging.DLQConfig{
			ExchangeName: events.DLXExchange,
			RoutingKey:   events.GetDLQKey(a.config.ResultsQueue),
		},
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	err = rabbit.Subscribe(
		events.GetDLQQueueName(a.config.ResultsQueue),
		events.GetDLQKey(a.config.ResultsQueue),
		events.DLXExchange,
		deadLetterListener.Handle,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

}
