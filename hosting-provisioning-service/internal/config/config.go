package config

import "time"

type Config struct {
	AMQP_URL         string
	ProvisionQueue   string
	ProvisioningTime time.Duration
}

func Load() *Config {
	return &Config{
		AMQP_URL:         "amqp://guest:guest@localhost:5672/",
		ProvisionQueue:   "provisioning_queue",
		ProvisioningTime: 10 * time.Second,
	}
}
