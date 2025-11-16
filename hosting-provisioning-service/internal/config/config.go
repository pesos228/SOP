package config

import "time"

type Config struct {
	AMQP_URL            string
	AMQP_HandlerTimeout time.Duration
	ProvisionQueue      string
	ProvisioningTime    time.Duration
}

func Load() *Config {
	return &Config{
		AMQP_URL:            "amqp://guest:guest@localhost:5672/",
		AMQP_HandlerTimeout: 10 * time.Second,
		ProvisionQueue:      "provisioning_queue",
		ProvisioningTime:    10 * time.Second,
	}
}
