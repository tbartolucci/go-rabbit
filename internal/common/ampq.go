package common

import (
	"github.com/streadway/amqp"
	"log"
)

type AmqpRabbitMQ struct {
	URL string
	Connection *amqp.Connection
	Channel *amqp.Channel
	Queue *amqp.Queue
}

func CreateAmpqRabbitMQ() AmqpRabbitMQ {
	return AmqpRabbitMQ{
		URL: "amqp://guest:guest@localhost:5672/",
	}
}

func (a *AmqpRabbitMQ) Close() {
	a.Channel.Close()
	a.Connection.Close()
}

func (a *AmqpRabbitMQ) GetChannel() *amqp.Channel{
	conn, err := amqp.Dial(a.URL)
	HandleError(err, "Con't connect to Rabbit")

	a.Connection = conn

	amqpChannel, err := conn.Channel()
	HandleError(err, "Can't create a amqpChannel")

	a.Channel = amqpChannel

	return a.Channel
}

func (a *AmqpRabbitMQ) GetQueue(queueName string) amqp.Queue {
	if a.Channel == nil {
		a.GetChannel()
	}
	queue, err := a.Channel.QueueDeclare(queueName, true, false, false, false, nil)
	HandleError(err, "Couldn't declar `add`  queue")
	a.Queue = &queue
	return queue
}

func (a *AmqpRabbitMQ) GetMessageChannel(queueName string) <-chan amqp.Delivery{
	if a.Channel == nil {
		a.GetChannel()
	}
	if a.Queue == nil {
		a.GetQueue(queueName)
	}

	err := a.Channel.Qos(1,0, false)
	HandleError(err, "Could not configure QoS")

	messageChannel, err := a.Channel.Consume(
		a.Queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	HandleError(err, "Could not register consumer")

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