package amqp

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strconv"
)

type RabbitMqService struct {
	url string
}

func NewRabbitMqService(user string, password string, host string, port int) *RabbitMqService {
	rabbitMqService := new(RabbitMqService)
	url := "amqp://" + user + ":" + password + "@" + host + ":" + strconv.Itoa(port) + "/"
	log.Info("amqp " + url)
	rabbitMqService.url = url
	return rabbitMqService
}

func (r *RabbitMqService) Publish(exchange string, routingKey string, body string) error {
	conn, err := amqp.Dial(r.url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare an exchange", err)
		return err
	}

	err = ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish a message", err)
		return err
	}

	log.Printf(" [x] Sent %s", body)
	return nil
}

func (r *RabbitMqService) Consume(exchange string, routingKey string, handler ConsumerHandler) {
	conn, err := amqp.Dial(r.url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	err = channel.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	queue, err := channel.QueueDeclare(
		"",
		true,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	log.Printf("Binding queue %s to exchange %s with routing key %s", queue.Name, exchange, routingKey)

	err = channel.QueueBind(
		queue.Name,
		routingKey,
		exchange,
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)
	go func() {
		for data := range messages {
			r.handleMessage(data, handler)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func (r *RabbitMqService) handleMessage(data amqp.Delivery, handler ConsumerHandler) {
	message := string(data.Body)
	err := handler.Handle(message)

	var args string
	if err != nil {
		args = "Failed"
	} else {
		args = "Success"
	}

	log.Info(args + " handle message '" + message + "'")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
