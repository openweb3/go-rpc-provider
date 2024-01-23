package providers

import (
	"context"
	"io"
	"time"

	"github.com/mcuadros/go-defaults"
)

// Option for set retry and timeout options
// Note: user could overwrite RequestTimeout when CallContext with timeout context or cancel context
type Option struct {
	// retry
	RetryCount    int
	RetryInterval time.Duration `default:"1s"`

	RequestTimeout time.Duration `default:"30s"`

	MaxConnectionPerHost int
	Logger               io.Writer
	CircuitBreaker       CircuitBreaker
}

func (o *Option) WithRetry(retryCount int, retryInterval time.Duration) *Option {
	o.RetryCount = retryCount
	o.RetryInterval = retryInterval
	return o
}

func (o *Option) WithTimout(requestTimeout time.Duration) *Option {
	o.RequestTimeout = requestTimeout
	return o
}

func (o *Option) WithMaxConnectionPerHost(maxConnectionPerHost int) *Option {
	o.MaxConnectionPerHost = maxConnectionPerHost
	return o
}

func (o *Option) WithLooger(w io.Writer) *Option {
	o.Logger = w
	return o
}

func (o *Option) WithCircuitBreaker(circuitBreakerOption DefaultCircuitBreakerOption) *Option {
	o.CircuitBreaker = NewDefaultCircuitBreaker(circuitBreakerOption)
	return o
}

// NewProviderWithOption returns a new MiddlewareProvider with hook handlers build according to options
// Note: user could overwrite RequestTimeout when CallContext with timeout context or cancel context
func NewProviderWithOption(rawurl string, option Option) (*MiddlewarableProvider, error) {
	p, err := NewBaseProvider(context.Background(), rawurl, option.MaxConnectionPerHost)
	if err != nil {
		return nil, err
	}

	defaults.SetDefaults(&option)
	p = wrapProvider(p, option)
	return p, nil
}

// wrapProvider wrap provider accroding to option
func wrapProvider(p *MiddlewarableProvider, option Option) *MiddlewarableProvider {
	if option.CircuitBreaker != nil {
		p = NewCircuitBreakerProvider(p, option.CircuitBreaker)
	}
	p = NewTimeoutableProvider(p, option.RequestTimeout)
	p = NewRetriableProvider(p, option.RetryCount, option.RetryInterval)
	p = NewLoggerProvider(p, option.Logger)
	return p
}
