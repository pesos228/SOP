package app

import (
	events "hosting-events-contract"
	"hosting-kit/messaging"
	"hosting-provisioning-service/internal/listeners"
	"hosting-provisioning-service/internal/service"
	"log"
	"sync"
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

func (a *App) runConsumers(rabbit *messaging.MessageManager, service service.ProvisioningService) {
	provisioningListener := listeners.NewProvisioningListener(a.config.ProvisioningTime, service)

	err := rabbit.Subscribe(
		a.config.ProvisionQueue,
		events.ProvisionRequestKey,
		events.CommandsExchange,
		provisioningListener.Handle,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
}
