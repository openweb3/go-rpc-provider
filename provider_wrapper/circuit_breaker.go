package providers

import (
	"errors"
	"time"

	"github.com/openweb3/go-rpc-provider/utils"
)

type ICircuitBreaker interface {
	Do(handler func() error) error
	State() BreakerState
}

type BreakerState int

const (
	BREAKER_CLOSED BreakerState = iota
	BREAKER_HALF_OPEN
	BREAKER_OPEN
)

var ErrCircuitOpen = errors.New("circuit breaked")
var ErrUnknownCircuitState = errors.New("unknown circuit state")

type CircuitBreaker struct {
	MaxFail        int
	FailTimeWindow time.Duration // continuous fail maxFail times between failTimeWindow, close -> open
	OpenColdTime   time.Duration // after openColdTime, open -> halfopen
	failHistory    []time.Time
	lastState      BreakerState // the state changed when Do
}

func (c *CircuitBreaker) Do(handler func() error) error {
	switch c.State() {
	case BREAKER_CLOSED:
		err := handler()
		if !utils.IsRPCJSONError(err) {
			c.failHistory = append(c.failHistory, time.Now())
		}

		isReached, maxfailUsedTime := c.maxfailUsedTime()

		if !isReached || maxfailUsedTime > c.FailTimeWindow {
			c.lastState = BREAKER_CLOSED
			return err
		}

		if c.sinceLastFail() < c.OpenColdTime {
			c.lastState = BREAKER_OPEN
		} else {
			c.lastState = BREAKER_HALF_OPEN
		}

		return err

	case BREAKER_HALF_OPEN:
		c.failHistory = []time.Time{}

		err := handler()
		if err == nil || utils.IsRPCJSONError(err) {
			c.lastState = BREAKER_CLOSED
		} else {
			c.failHistory = append(c.failHistory, time.Now())
			c.lastState = BREAKER_OPEN
		}
		return err

	case BREAKER_OPEN:
		return ErrCircuitOpen

	default:
		return ErrUnknownCircuitState
	}
}

func (c *CircuitBreaker) State() BreakerState {
	if c.lastState == BREAKER_OPEN && c.sinceLastFail() > c.OpenColdTime {
		return BREAKER_HALF_OPEN
	}

	return c.lastState
}

// 1st return means if reached max fail.
func (c *CircuitBreaker) maxfailUsedTime() (bool, time.Duration) {
	failLen := len(c.failHistory)
	if failLen < c.MaxFail {
		return false, 0
	}

	return true, time.Since(c.failHistory[failLen-c.MaxFail]) - c.sinceLastFail()
}

func (c *CircuitBreaker) sinceLastFail() time.Duration {
	lastFail := c.failHistory[len(c.failHistory)-1]
	return time.Since(lastFail)
}
