package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/josemarjobs/sensors_app/dto"
	"github.com/josemarjobs/sensors_app/qutils"
	"github.com/streadway/amqp"
)

const url = "amqp://guest:guest@localhost:5672"

type QueueListener struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources map[string]<-chan amqp.Delivery
	ea      *EventAggregator
}

func NewQueueListener(ea *EventAggregator) *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
		ea:      ea,
	}

	ql.conn, ql.ch = qutils.GetChannel(url)
	return &ql
}

func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("", ql.ch, true)
	ql.ch.QueueBind(
		q.Name,       // name string
		"",           // key string
		"amq.fanout", // exchange string
		false,        // noWait bool
		nil)

	msgs, _ := ql.ch.Consume(q.Name, "", false, false, false, false, nil)

	ql.DiscoverSensors()

	log.Println("listening for new sources")
	for msg := range msgs {
		log.Println("new source discovered")
		ql.ea.PublishEvent("DataSourceDiscovered", string(msg.Body))

		sourceChan, _ := ql.ch.Consume(
			string(msg.Body), "", true, false, false, false, nil)
		if ql.sources[string(msg.Body)] == nil {
			ql.sources[string(msg.Body)] = sourceChan

			go ql.AddListener(sourceChan)
		}
	}
}

func (ql *QueueListener) DiscoverSensors() {
	ql.ch.ExchangeDeclare(
		qutils.SensorDiscoveryExchange, // name
		"fanout",                       // kind string
		false,                          // durable bool
		false,                          // autoDelete bool
		false,                          // internal bool
		false,                          // noWait bool
		nil,                            // args amqp.Table
	)
	ql.ch.Publish(
		qutils.SensorDiscoveryExchange, // name string
		"",                // key string
		false,             // mandatory bool
		false,             // immediate bool
		amqp.Publishing{}, // msg amqp.Publishing
	)
}
func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sd := new(dto.SensorMessage)
		d.Decode(sd)

		fmt.Printf("Received Message: %v\n", sd)

		ed := EventData{
			Name:      sd.Name,
			Timestamp: sd.Timestamp,
			Value:     sd.Value,
		}
		ql.ea.PublishEvent("MessageReceived_"+msg.RoutingKey, ed)
	}
}
