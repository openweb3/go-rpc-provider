package providers

import (
	"context"
	"time"

	"github.com/mcuadros/go-defaults"
)

// Option for set retry and timeout options
// Note: user could overwrite RequestTimeout when CallContext with timeout context or cancel context
type Option struct {
	// KeystorePath string
	// retry
	RetryCount    int           `default:"3"`
	RetryInterval time.Duration `default:"1s"`
	// timeout of request
	RequestTimeout time.Duration `default:"3s"`
	// Maximum number of connections may be established. The default value is 512.
	MaxConnectionPerHost int
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
	p = NewTimeoutableProvider(p, option.RequestTimeout)
	p = NewRetriableProvider(p, option.RetryCount, option.RetryInterval)
	return p
}
