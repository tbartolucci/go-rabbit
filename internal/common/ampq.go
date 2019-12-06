package common

import (
	"log"

	"github.com/streadway/amqp"
)

type AmqpRabbitMQ struct {
	URL        string
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      *amqp.Queue
}

func CreateAmpqRabbitMQ(url string) AmqpRabbitMQ {
	return AmqpRabbitMQ{
		URL: url,
	}
}

func (a *AmqpRabbitMQ) Close() {
	a.Channel.Close()
	a.Connection.Close()
}

func (a *AmqpRabbitMQ) handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (a *AmqpRabbitMQ) GetChannel() *amqp.Channel {
	conn, err := amqp.Dial(a.URL)
	a.handleError(err, "Con't connect to Rabbit")

	a.Connection = conn

	amqpChannel, err := conn.Channel()
	a.handleError(err, "Can't create a amqpChannel")

	a.Channel = amqpChannel

	return a.Channel
}

func (a *AmqpRabbitMQ) GetQueue(queueName string) amqp.Queue {
	if a.Channel == nil {
		a.GetChannel()
	}
	queue, err := a.Channel.QueueDeclare(queueName, true, false, false, false, nil)
	a.handleError(err, "Couldn't declar `add`  queue")
	a.Queue = &queue
	return queue
}

func (a *AmqpRabbitMQ) GetMessageChannel(queueName string) <-chan amqp.Delivery {
	if a.Channel == nil {
		a.GetChannel()
	}
	if a.Queue == nil {
		a.GetQueue(queueName)
	}

	err := a.Channel.Qos(1, 0, false)
	a.handleError(err, "Could not configure QoS")

	messageChannel, err := a.Channel.Consume(
		a.Queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	a.handleError(err, "Could not register consumer")

	return messageChannel
}

func (a *AmqpRabbitMQ) Publish(queueName string, body []byte) {
	err := a.Channel.Publish("", a.Queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	log.Printf("%s", body)

}
