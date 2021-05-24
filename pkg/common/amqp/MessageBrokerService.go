package amqp

type ConsumerHandler interface {
	Handle(message string) error
}

type MessageBroker interface {
	Publish(exchange string, routingKey string, body string)
	Consume(exchange string, routingKey string, handler ConsumerHandler)
}
