package testx

import (
	"context"
	"errors"
	"time"
)

var errTestDeadlineExceeded = errors.New("deadline exceeded")

func Context(t testT) context.Context {
	deadline, ok := t.Deadline()
	if !ok {
		ctx, stop := context.WithCancelCause(context.Background())
		t.Cleanup(func() {
			stop(errTestDeadlineExceeded)
		})
		return ctx
	} else {
		ctx, stop := context.WithDeadlineCause(context.Background(), deadline, errTestDeadlineExceeded)
		t.Cleanup(stop)
		return ctx
	}
}

type testT interface {
	Deadline() (deadline time.Time, ok bool)
	Cleanup(f func())
}
