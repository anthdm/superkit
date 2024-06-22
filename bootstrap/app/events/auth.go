package events

import (
	"AABBCCDD/plugins/auth"
	"context"
	"encoding/json"
	"fmt"
)

// Event handlers
func OnUserSignup(ctx context.Context, event any) {
	userWithToken, ok := event.(auth.UserWithVerificationToken)
	if !ok {
		return
	}
	b, _ := json.MarshalIndent(userWithToken, "   ", "    ")
	fmt.Println(string(b))
}

func OnResendVerificationToken(ctx context.Context, event any) {
	userWithToken, ok := event.(auth.UserWithVerificationToken)
	if !ok {
		return
	}
	b, _ := json.MarshalIndent(userWithToken, "   ", "    ")
	fmt.Println(string(b))
}
