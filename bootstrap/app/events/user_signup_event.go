package events

import (
	"AABBCCDD/plugins/auth"
	"context"
	"fmt"
)

// Event handlers
func HandleUserSignup(ctx context.Context, event any) {
	userWithToken, ok := event.(auth.UserWithVerificationToken)
	if !ok {
		return
	}
	fmt.Printf("user signup: %v\n", userWithToken)
}
