package event

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func TestEventSubscribeEmit(t *testing.T) {
	var (
		expect = 1
		wg     = sync.WaitGroup{}
	)
	wg.Add(1)
	Subscribe("foo.bar", func(_ context.Context, event any) {
		value, ok := event.(int)
		if !ok {
			t.Errorf("expected int got %v", reflect.TypeOf(event))
		}
		if value != 1 {
			t.Errorf("expected %d got %d", expect, value)
		}
		wg.Done()
	})
	Emit("foo.bar", expect)
	wg.Wait()
}

func TestUnsubscribe(t *testing.T) {
	sub := Subscribe("foo.bar", func(_ context.Context, _ any) {})
	Unsubscribe(sub)
	if _, ok := stream.subs["foo.bar"]; ok {
		t.Errorf("expected topic foo.bar to be deleted")
	}
}
