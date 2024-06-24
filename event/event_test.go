package event

import (
	"context"
	"reflect"
	"testing"
)

func TestEventSubscribeEmit(t *testing.T) {
	expect := 1
	ctx, cancel := context.WithCancel(context.Background())
	Subscribe("foo.a", func(_ context.Context, event any) {
		defer cancel()
		value, ok := event.(int)
		if !ok {
			t.Errorf("expected int got %v", reflect.TypeOf(event))
		}
		if value != 1 {
			t.Errorf("expected %d got %d", expect, value)
		}
	})
	Emit("foo.a", expect)
	<-ctx.Done()
}

func TestUnsubscribe(t *testing.T) {
	sub := Subscribe("foo.b", func(_ context.Context, _ any) {})
	Unsubscribe(sub)
	if _, ok := stream.subs["foo.b"]; ok {
		t.Errorf("expected topic foo.bar to be deleted")
	}
}
