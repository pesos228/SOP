module hosting-provisioning-service

go 1.24.4

require github.com/wagslane/go-rabbitmq v0.15.0

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	hosting-events-contract v0.0.0-00010101000000-000000000000
)

replace hosting-events-contract => ../hosting-events-contract
