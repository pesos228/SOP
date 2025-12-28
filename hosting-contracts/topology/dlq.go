package topology

const (
	DLXExchange         = "hosting.dlx"
	DLQRoutingKeyPrefix = "dlq."
	DLQQueueSuffix      = ".dlq"
)

func GetDLQKey(queue string) string {
	return DLQRoutingKeyPrefix + queue
}

func GetDLQQueueName(originalQueue string) string {
	return originalQueue + DLQQueueSuffix
}
