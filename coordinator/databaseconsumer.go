package coordinator

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/josemarjobs/sensors_app/dto"
	"github.com/josemarjobs/sensors_app/qutils"
	"github.com/streadway/amqp"
)

const (
	maxRate = 5 * time.Second
)

type DatabaseConsumer struct {
	er      EventRaiser
	conn    *amqp.Connection
	ch      *amqp.Channel
	queue   *amqp.Queue
	sources []string
}

func NewDatabaseConsumer(er EventRaiser) *DatabaseConsumer {
	dc := DatabaseConsumer{er: er}

	dc.conn, dc.ch = qutils.GetChannel(url)
	dc.queue = qutils.GetQueue(qutils.PersistenceReadingsQueue, dc.ch, false)

	dc.er.AddListener("DataSourceDiscovered", func(eventData interface{}) {
		dc.SubscribeToDataEvent(eventData.(string))
	})
	return &dc
}

func (dc *DatabaseConsumer) SubscribeToDataEvent(eventName string) {
	for _, source := range dc.sources {
		if source == eventName {
			return
		}
	}

	dc.er.AddListener("MessageReceived_"+eventName,
		func() func(interface{}) {
			prevTime := time.Unix(0, 0)
			buf := new(bytes.Buffer)
			return func(eventData interface{}) {
				ed := eventData.(EventData)
				if time.Since(prevTime) > maxRate {
					prevTime = time.Now()

					sm := dto.SensorMessage{
						Name:      ed.Name,
						Value:     ed.Value,
						Timestamp: ed.Timestamp,
					}
					buf.Reset()

					gob.NewEncoder(buf).Encode(sm)

					msg := amqp.Publishing{Body: buf.Bytes()}

					dc.ch.Publish(
						"",
						qutils.PersistenceReadingsQueue,
						false,
						false,
						msg,
					)
				}
			}
		}())
}
