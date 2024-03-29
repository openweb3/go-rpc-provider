package providers

import (
	"errors"
	"math"
	"sync"
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

type DefaultCircuitBreakerOption struct {
	MaxFail        int           `default:"5"`
	FailTimeWindow time.Duration `default:"10s"` // continuous fail maxFail times between failTimeWindow, close -> open
	OpenColdTime   time.Duration `default:"10s"` // after openColdTime, open -> halfopen
}

type DefaultCircuitBreaker struct {
	DefaultCircuitBreakerOption
	failHistory []time.Time
	lastState   BreakerState // the state changed when Do
	sync.Mutex
}

func NewDefaultCircuitBreaker(option ...DefaultCircuitBreakerOption) *DefaultCircuitBreaker {
	if len(option) == 0 {
		option = append(option, DefaultCircuitBreakerOption{})
	}
	defaults.SetDefaults(&option[0])

	return &DefaultCircuitBreaker{
		DefaultCircuitBreakerOption: option[0],
	}
}

func (c *DefaultCircuitBreaker) Do(handler func() error) error {
	switch c.State() {
	case BREAKER_CLOSED:
		err := handler()

		c.Lock()
		defer c.Unlock()

		if err == nil || utils.IsRPCJSONError(err) {
			c.failHistory = []time.Time{}
			c.lastState = BREAKER_CLOSED
			return err
		} else {
			c.failHistory = append(c.failHistory, time.Now())
			if !c.isMaxfailAcheived() {
				return err
			}

			c.lastState = BREAKER_OPEN
		}

		return err

	case BREAKER_HALF_OPEN:
		err := handler()

		c.Lock()
		defer c.Unlock()

		c.failHistory = []time.Time{}
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
	c.Lock()
	defer c.Unlock()

	if c.lastState == BREAKER_OPEN && c.sinceLastFail() > c.OpenColdTime {
		return BREAKER_HALF_OPEN
	}

	return c.lastState
}

func (c *DefaultCircuitBreaker) isMaxfailAcheived() bool {
	isReached, maxfailUsedTime := c.maxfailUsedTime()

	if !isReached || maxfailUsedTime > c.FailTimeWindow {
		return false
	}
	return true
}

// 1st return means if reached max fail.
func (c *DefaultCircuitBreaker) maxfailUsedTime() (bool, time.Duration) {
	failLen := len(c.failHistory)
	if failLen < c.MaxFail {
		return false, 0
	}

	return true, time.Since(c.failHistory[failLen-c.MaxFail]) - c.sinceLastFail()
}

// note: return max int64 when failHistory is empty becase that means no failure before.
func (c *DefaultCircuitBreaker) sinceLastFail() time.Duration {
	if len(c.failHistory) == 0 {
		return math.MaxInt64
	}
	lastFail := c.failHistory[len(c.failHistory)-1]
	return time.Since(lastFail)
}
