package qutils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const SensorDiscoveryExchange = "SensorDiscovery"
const PersistenceReadingsQueue = "PersistenceReadingsQueue"

func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to establish connection to msg broker")
	ch, err := conn.Channel()
	failOnError(err, "Failed to get channel for connection")

	return conn, ch
}

func GetQueue(name string, ch *amqp.Channel, autoDelete bool) *amqp.Queue {
	q, err := ch.QueueDeclare(name, false, autoDelete, false, false, nil)
	failOnError(err, "Error declaring the queue")
	return &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
