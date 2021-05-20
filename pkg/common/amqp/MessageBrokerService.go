package amqp

type MessageBroker interface {
	Publish(exchange string, routingKey string, body string)
}
