package app

import (
	"AABBCCDD/app/events"

	"github.com/anthdm/gothkit/event"
)

// Events are functions that are handled in separate goroutines.
// They are the perfect fit for offloading work in your handlers
// that otherwise would take up response time.
// - sending email
// - sending notifications (Slack, Telegram, Discord)
// - analytics..

// Register your events here.
func RegisterEvents() {
	event.Subscribe("foo.bar", events.HandleFooEvent)
}
