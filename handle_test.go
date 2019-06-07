package interrupt_test

import (
	"context"
	"os"
	"testing"

	"github.com/torpidtangerine/interrupt"
)

func TestWithInterrupt(t *testing.T) {
	expVal := "val"
	ctx := context.WithValue(context.Background(), "key", expVal)
	ctx, cancel := interrupt.WithInterrupt(ctx)
	defer cancel()

	ctxVal := ctx.Value("key")
	if ctxVal != expVal {
		t.Errorf("expected '%v' received '%v'", expVal, ctxVal)
	}

	if ctx.Err() == context.Canceled {
		t.Errorf("expected ctx to not be canceled")
	}
	cancel()
	if ctx.Err() != context.Canceled {
		t.Errorf("expected ctx to be canceled")
	}
}

func TestWithInterruptSignal(t *testing.T) {
	var signalChan chan<- os.Signal
	stopChan := make(chan struct{})

	notifier := interrupt.NewNotifier(
		func(c chan<- os.Signal, sig ...os.Signal) {
			signalChan = c
			if len(sig) != 1 {
				t.Errorf("expected 1 signal")
				return
			}

			if sig[0] != os.Interrupt {
				t.Errorf("expected os.Interrupt")
			}
		},
		func(c chan<- os.Signal) {
			signalChan = nil
			stopChan <- struct{}{}
		},
	)

	ctx, cancel := notifier.Background()
	defer cancel()

	if ctx.Err() == context.Canceled {
		t.Errorf("expected not to be canceled")
	}
	signalChan <- os.Interrupt

	// handles race condition - context cancelation is not immediate, but is always before stop
	<-stopChan

	if ctx.Err() != context.Canceled {
		t.Errorf("expected to be canceled")
	}
}
