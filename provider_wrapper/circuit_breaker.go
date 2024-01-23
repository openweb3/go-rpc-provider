package providers

import (
	"errors"
	"time"

	"github.com/mcuadros/go-defaults"
	"github.com/openweb3/go-rpc-provider/utils"
)

type CircuitBreaker interface {
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

type DefaultCircuitBreaderOption struct {
	MaxFail        int           `default:"5"`
	FailTimeWindow time.Duration `default:"10s"` // continuous fail maxFail times between failTimeWindow, close -> open
	OpenColdTime   time.Duration `default:"10s"` // after openColdTime, open -> halfopen
}

type DefaultCircuitBreaker struct {
	DefaultCircuitBreaderOption
	failHistory []time.Time
	lastState   BreakerState // the state changed when Do
}

func NewDefaultCircuitBreaker(option ...DefaultCircuitBreaderOption) *DefaultCircuitBreaker {
	if len(option) == 0 {
		option = []DefaultCircuitBreaderOption{}
	}
	defaults.SetDefaults(&option[0])

	return &DefaultCircuitBreaker{
		DefaultCircuitBreaderOption: option[0],
	}
}

func (c *DefaultCircuitBreaker) Do(handler func() error) error {
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

func (c *DefaultCircuitBreaker) State() BreakerState {
	if c.lastState == BREAKER_OPEN && c.sinceLastFail() > c.OpenColdTime {
		return BREAKER_HALF_OPEN
	}

	return c.lastState
}

// 1st return means if reached max fail.
func (c *DefaultCircuitBreaker) maxfailUsedTime() (bool, time.Duration) {
	failLen := len(c.failHistory)
	if failLen < c.MaxFail {
		return false, 0
	}

	return true, time.Since(c.failHistory[failLen-c.MaxFail]) - c.sinceLastFail()
}

func (c *DefaultCircuitBreaker) sinceLastFail() time.Duration {
	lastFail := c.failHistory[len(c.failHistory)-1]
	return time.Since(lastFail)
}
