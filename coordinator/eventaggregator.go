package coordinator

import "time"

type Callback func(interface{})

type EventRaiser interface {
	AddListener(eventName string, f Callback)
}

type EventAggregator struct {
	listeners map[string][]Callback
}

func NewEventAggregator() *EventAggregator {
	return &EventAggregator{
		listeners: make(map[string][]Callback),
	}
}

func (ea *EventAggregator) AddListener(name string, f Callback) {
	ea.listeners[name] = append(ea.listeners[name], f)
}

func (ea *EventAggregator) PublishEvent(name string, eventData interface{}) {
	if ea.listeners[name] != nil {
		for _, cb := range ea.listeners[name] {
			cb(eventData)
		}
	}
}

type EventData struct {
	Name      string
	Value     float64
	Timestamp time.Time
}
