package messaging

type ExchangeType string

const (
	ExchangeDirect  ExchangeType = "direct"
	ExchangeFanout  ExchangeType = "fanout"
	ExchangeTopic   ExchangeType = "topic"
	ExchangeHeaders ExchangeType = "headers"
)

type ExchangeConfig struct {
	Name string
	Type ExchangeType
}
