package process

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetProcessInfoChan(t *testing.T) {
	t.Run("should close channel on context cancel", func(t *testing.T) {
		// given
		ctx, cancel := context.WithCancel(context.Background())
		interval := 100 * time.Millisecond
		infoChan := GetProcessInfoChan(ctx, interval)

		// when
		cancel()

		// then
		select {
		case _, ok := <-infoChan:
			assert.False(t, ok, "channel should be closed")
		case <-time.After(1 * time.Second):
			assert.Fail(t, "channel not closed after 1 second")
		}
	})

	t.Run("should receive process info", func(t *testing.T) {
		// given
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		interval := 100 * time.Millisecond

		// when
		infoChan := GetProcessInfoChan(ctx, interval)

		// then
		var receivedProcess bool
	loop:
		for {
			select {
			case info, ok := <-infoChan:
				assert.True(t, ok, "should be able to receive on channel")
				if info.PID != 0 {
					receivedProcess = true
					break loop
				}
			case <-ctx.Done():
				break loop
			}
		}

		assert.True(t, receivedProcess, "did not receive any process info with non-zero PID before timeout")
	})
}
