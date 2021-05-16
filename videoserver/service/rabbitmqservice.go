package service

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"strconv"
	"strings"
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
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	//body := bodyFrom(os.Args)
	err = ch.Publish(
		exchange,          // exchange
		routingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
	//
	//channel, err := r.connection.Channel()
	//failOnError(err, "Failed to open a channel")
	//defer channel.Close()
	//err = channel.ExchangeDeclare(
	//	exchange, // name
	//	"topic",      // type
	//	true,         // durable
	//	false,        // auto-deleted
	//	false,        // internal
	//	false,        // no-wait
	//	nil,          // arguments
	//)
	//failOnError(err, "Failed to declare an exchange")

	//err = r.channel.Publish(
	//	"logs_topic",          // exchange
	//	severityFrom(os.Args), // routing key
	//	false,                 // mandatory
	//	false,                 // immediate
	//	amqp.Publishing{
	//		ContentType: "text/plain",
	//		Body:        []byte(body),
	//	})
	//err = channel.Publish(
	//	exchange,          // exchange
	//	routingKey, // routing key
	//	false,                 // mandatory
	//	false,                 // immediate
	//	amqp.Publishing{
	//		ContentType: "text/plain",
	//		Body:        []byte(body),
	//	})
	//
	//failOnError(err, "Failed to publish a message")
	//
	//log.Printf(" [x] Sent %s", body)
}

func (r *RabbitMqService) publishEvent() {

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}
