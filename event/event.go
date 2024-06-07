package event

import (
	"context"
	"slices"
	"sync"
	"time"
)

// HandlerFunc is the function being called when receiving an event.
type HandlerFunc func(context.Context, any)

// Emit and event to the given topic
func Emit(topic string, event any) {
	stream.emit(topic, event)
}

// Subscribe a HandlerFunc to the given topic.
// A Subscription is being returned that can be used
// to unsubscribe from the topic.
func Subscribe(topic string, h HandlerFunc) Subscription {
	return stream.subscribe(topic, h)
}

// Unsubscribe unsubribes the given Subscription from its topic.
func Unsubscribe(sub Subscription) {
	stream.unsubscribe(sub)
}

// Stop stops the event stream, cleaning up its resources.
func Stop() {
	stream.stop()
}

var stream *eventStream

type event struct {
	topic   string
	message any
}

// Subscription represents a handler subscribed to a specific topic.
type Subscription struct {
	Topic     string
	CreatedAt int64
	Fn        HandlerFunc
}

type eventStream struct {
	mu      sync.RWMutex
	subs    map[string][]Subscription
	eventch chan event
	quitch  chan struct{}
}

func newStream() *eventStream {
	e := &eventStream{
		subs:    make(map[string][]Subscription),
		eventch: make(chan event, 128),
		quitch:  make(chan struct{}),
	}
	go e.start()
	return e
}

func (e *eventStream) start() {
	ctx := context.Background()
	for {
		select {
		case <-e.quitch:
			return
		case evt := <-e.eventch:
			if handlers, ok := e.subs[evt.topic]; ok {
				for _, sub := range handlers {
					go sub.Fn(ctx, evt.message)
				}
			}
		}
	}
}

func (e *eventStream) stop() {
	e.quitch <- struct{}{}
}

func (e *eventStream) emit(topic string, v any) {
	e.eventch <- event{
		topic:   topic,
		message: v,
	}
}

func (e *eventStream) subscribe(topic string, h HandlerFunc) Subscription {
	e.mu.RLock()
	defer e.mu.RUnlock()

	sub := Subscription{
		CreatedAt: time.Now().UnixNano(),
		Topic:     topic,
		Fn:        h,
	}

	if _, ok := e.subs[topic]; !ok {
		e.subs[topic] = []Subscription{}
	}

	e.subs[topic] = append(e.subs[topic], sub)

	return sub
}

func (e *eventStream) unsubscribe(sub Subscription) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if _, ok := e.subs[sub.Topic]; ok {
		e.subs[sub.Topic] = slices.DeleteFunc(e.subs[sub.Topic], func(e Subscription) bool {
			return sub.CreatedAt == e.CreatedAt
		})
	}
	if len(e.subs[sub.Topic]) == 0 {
		delete(e.subs, sub.Topic)
	}
}

func init() {
	stream = newStream()
}
