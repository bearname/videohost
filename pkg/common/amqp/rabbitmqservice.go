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

func (r *RabbitMqService) Publish(exchange string, routingKey string, body string) {
	conn, err := amqp.Dial(r.url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
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
	failOnError(err, "Failed to declare an exchange")

	err = ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

func (r *RabbitMqService) publishEvent() {

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
