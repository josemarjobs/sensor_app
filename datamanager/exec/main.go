package main

import (
	"log"

	"bytes"
	"encoding/gob"

	"github.com/josemarjobs/sensors_app/datamanager"
	"github.com/josemarjobs/sensors_app/dto"
	"github.com/josemarjobs/sensors_app/qutils"
)

const url = "amqp://guest:guest@localhost:5672"

func main() {
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		qutils.PersistenceReadingsQueue, // name
		"",    // consumer
		false, // autoAck bool
		true,  // exclusive bool
		false, // noLocal bool
		false, // noWait bool
		nil,   // args amqp.Table
	)

	if err != nil {
		log.Fatalln("Failed to get access to messages")
	}

	for msg := range msgs {
		buf := bytes.NewReader(msg.Body)
		sd := new(dto.SensorMessage)
		gob.NewDecoder(buf).Decode(sd)

		err := datamanager.SaveReading(sd)
		if err != nil {
			log.Printf("failed to save reading from sensor %v. Error %v\n", sd.Name, err.Error())
		} else {
			msg.Ack(false)
		}
	}

}
