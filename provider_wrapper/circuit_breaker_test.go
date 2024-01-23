package providers

import (
	"errors"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestCircuitBreakerClose2Open(t *testing.T) {
	cb := DefaultCircuitBreaker{
		MaxFail:        3,
		FailTimeWindow: time.Second,
		OpenColdTime:   time.Millisecond * 500,
	}

	ErrTest := errors.New("test")
	// is closed
	assert.Equal(t, cb.State(), BREAKER_CLOSED)

	// closed -> open
	cb.Do(func() error { return ErrTest })
	cb.Do(func() error { return ErrTest })
	assert.Equal(t, cb.State(), BREAKER_CLOSED)

	time.Sleep(cb.FailTimeWindow)

	cb.Do(func() error { return ErrTest })
	assert.Equal(t, cb.State(), BREAKER_CLOSED)

	cb.Do(func() error { return ErrTest })
	cb.Do(func() error { return ErrTest })
	assert.Equal(t, cb.State(), BREAKER_OPEN)
}

func TestCircuitBreakerOpen2HalfOpen(t *testing.T) {
	cb := DefaultCircuitBreaker{
		MaxFail:        3,
		FailTimeWindow: time.Second,
		OpenColdTime:   time.Millisecond * 500,
		failHistory:    []time.Time{time.Now(), time.Now(), time.Now()},
		lastState:      BREAKER_OPEN,
	}

	// is open
	assert.Equal(t, cb.State(), BREAKER_OPEN)

	err := cb.Do(func() error { return errors.New("") })
	assert.Equal(t, err, ErrCircuitOpen)

	// open -> half open
	time.Sleep(cb.OpenColdTime)
	assert.Equal(t, cb.State(), BREAKER_HALF_OPEN)
}

func TestCircuitBreakerHalfOpen2Other(t *testing.T) {
	cb := DefaultCircuitBreaker{
		MaxFail:        3,
		FailTimeWindow: time.Second,
		OpenColdTime:   time.Millisecond * 500,
		lastState:      BREAKER_HALF_OPEN,
	}

	ErrTest := errors.New("test")

	// half open -> open
	err := cb.Do(func() error { return ErrTest })
	assert.Equal(t, err, ErrTest)
	assert.Equal(t, cb.State(), BREAKER_OPEN)

	// half open -> closed
	time.Sleep(cb.OpenColdTime)
	assert.Equal(t, cb.State(), BREAKER_HALF_OPEN)
	cb.Do(func() error { return nil })
	assert.Equal(t, cb.State(), BREAKER_CLOSED)
}
