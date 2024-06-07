package event

import (
	"context"
	"slices"
	"sync"
	"time"
)

type event struct {
	topic   string
	message any
}

type CreateUserEvent struct{}

func HandleCreateUserEvent(ctx context.Context, event CreateUserEvent) {}

// Handler is the function being called when receiving an event.
type Handler func(context.Context, any)

// Subscription represents a handler subscribed to a specific topic.
type Subscription struct {
	Topic     string
	CreatedAt int64
	Handler   Handler
}

// EventStream
type EventStream struct {
	mu      sync.RWMutex
	subs    map[string][]Subscription
	eventch chan event
	quitch  chan struct{}
}

// New return a new EventStream
func New() *EventStream {
	e := &EventStream{
		subs:    make(map[string][]Subscription),
		eventch: make(chan event, 128),
		quitch:  make(chan struct{}),
	}
	go e.start()
	return e
}

func (e *EventStream) start() {
	ctx := context.Background()
	for {
		select {
		case <-e.quitch:
			return
		case evt := <-e.eventch:
			if handlers, ok := e.subs[evt.topic]; ok {
				for _, sub := range handlers {
					go sub.Handler(ctx, evt.message)
				}
			}
		}
	}
}

// Stop stops the EventStream
func (e *EventStream) Stop() {
	e.quitch <- struct{}{}
}

// Emit an event by specifying a topic and an arbitrary data type
func (e *EventStream) Emit(topic string, v any) {
	e.eventch <- event{
		topic:   topic,
		message: v,
	}
}

// Subscribe subscribes a handler to the given topic
func (e *EventStream) Subscribe(topic string, h Handler) Subscription {
	e.mu.RLock()
	defer e.mu.RUnlock()

	sub := Subscription{
		CreatedAt: time.Now().UnixNano(),
		Topic:     topic,
		Handler:   h,
	}

	if _, ok := e.subs[topic]; !ok {
		e.subs[topic] = []Subscription{}
	}

	e.subs[topic] = append(e.subs[topic], sub)

	return sub
}

// Unsubscribe unsubscribes the given Subscription
func (e *EventStream) Unsubscribe(sub Subscription) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if _, ok := e.subs[sub.Topic]; ok {
		e.subs[sub.Topic] = slices.DeleteFunc(e.subs[sub.Topic], func(e Subscription) bool {
			return sub.CreatedAt == e.CreatedAt
		})
	}
}
