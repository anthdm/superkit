package app

import (
	"AABBCCDD/app/events"
	"AABBCCDD/plugins/auth"

	"github.com/anthdm/superkit/event"
)

// Events are functions that are handled in separate goroutines.
// They are the perfect fit for offloading work in your handlers
// that otherwise would take up response time.
// - sending email
// - sending notifications (Slack, Telegram, Discord)
// - analytics..

// Register your events here.
func RegisterEvents() {
	event.Subscribe(auth.UserSignupEvent, events.OnUserSignup)
	event.Subscribe(auth.ResendVerificationEvent, events.OnResendVerificationToken)
}
