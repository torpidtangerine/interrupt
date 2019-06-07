package interrupt

import (
	"context"
	"os"
	"os/signal"
)

// Notifier produces contexts that handle signals
type Notifier struct {
	notify func(chan<- os.Signal, ...os.Signal)
	stop   func(chan<- os.Signal)
}

// NewNotifier creates a Notifier
func NewNotifier(
	notify func(chan<- os.Signal, ...os.Signal),
	stop func(chan<- os.Signal),
) *Notifier {
	return &Notifier{
		notify: notify,
		stop:   stop,
	}
}

// WithInterrupt returns a new context and a cancel func which handles os.Interrupt
func (n *Notifier) WithInterrupt(ctx context.Context) (context.Context, context.CancelFunc) {
	ctxInner, cancel := context.WithCancel(ctx)

	signals := make(chan os.Signal)
	n.notify(signals, os.Interrupt)

	go func() {
		select {
		case <-signals:
			cancel()
		case <-ctxInner.Done():
		}

		n.stop(signals)
	}()

	return ctxInner, cancel
}

// Background returns a background context and a cancel func which handles os.Interrupt
func (n *Notifier) Background() (context.Context, context.CancelFunc) {
	return n.WithInterrupt(context.Background())
}

// Default is the default notifier
var Default = NewNotifier(signal.Notify, signal.Stop)

// WithInterrupt returns a new context and a cancel func which handles os.Interrupt
var WithInterrupt = Default.WithInterrupt

// Background returns a background context and a cancel func which handles os.Interrupt
var Background = Default.Background
