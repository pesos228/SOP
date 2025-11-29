package listeners

import (
	"context"
	"log"
)

type DeadLetterListener struct {
}

func NewDeadLetterListener() *DeadLetterListener {
	return &DeadLetterListener{}
}

func (l *DeadLetterListener) Handle(ctx context.Context, body []byte, routingKey string) error {
	log.Printf("Message dropped! RoutingKey: %s", routingKey)
	log.Printf("Body: %s", string(body))

	return nil
}
