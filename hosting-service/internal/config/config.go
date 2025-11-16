package config

import "time"

type Config struct {
	HTTP_Port           string
	DB_DSN              string
	AMQP_URL            string
	AMQP_HandlerTimeout time.Duration
	ResultsQueue        string
}

func Load() *Config {
	return &Config{
		HTTP_Port:           ":8080",
		DB_DSN:              "postgres://postgres:vladick@localhost:5432/sop?search_path=public",
		AMQP_URL:            "amqp://guest:guest@localhost:5672/",
		AMQP_HandlerTimeout: time.Second * 10,
		ResultsQueue:        "api_events_queue",
	}
}
