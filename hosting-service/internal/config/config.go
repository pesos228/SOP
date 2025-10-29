package config

type Config struct {
	HTTP_Port    string
	DB_DSN       string
	AMQP_URL     string
	ResultsQueue string
}

func Load() *Config {
	return &Config{
		HTTP_Port:    ":8080",
		DB_DSN:       "postgres://postgres:vladick@localhost:5432/sop?search_path=public",
		AMQP_URL:     "amqp://guest:guest@localhost:5672/",
		ResultsQueue: "api_events_queue",
	}
}
