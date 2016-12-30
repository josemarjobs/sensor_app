package coordinator

import "time"

type EventAggregator struct {
	listeners map[string][]Callback
}

type Callback func(EventData)

func NewEventAggregator() *EventAggregator {
	return &EventAggregator{
		listeners: make(map[string][]Callback),
	}
}

func (ea *EventAggregator) AddListener(name string, f Callback) {
	ea.listeners[name] = append(ea.listeners[name], f)
}

func (ea *EventAggregator) PublishEvent(name string, eventData EventData) {
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
