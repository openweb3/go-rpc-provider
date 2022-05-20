package providers

import (
	"context"
	"time"

	"github.com/mcuadros/go-defaults"
)

// Option for set retry and timeout options
type Option struct {
	// KeystorePath string
	// retry
	RetryCount    int           `default:"3"`
	RetryInterval time.Duration `default:"1s"`
	// timeout of request
	RequestTimeout time.Duration `default:"3s"`
	// Maximum number of connections may be established. The default value is 512.
	MaxConnectionNum int
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

func NewProviderWithOption(rawurl string, option *Option) (*MiddlewarableProvider, error) {
	maxConn := 0
	if option != nil {
		maxConn = option.MaxConnectionNum
	}

	p, err := NewBaseProvider(context.Background(), rawurl, maxConn)
	if err != nil {
		return nil, err
	}

	if option == nil {
		option = &Option{}
	}

	defaults.SetDefaults(option)
	p = wrapProvider(p, option)
	return p, nil
}

// wrapProvider wrap provider accroding to option
func wrapProvider(p *MiddlewarableProvider, option *Option) *MiddlewarableProvider {
	if option == nil {
		return p
	}

	p = NewTimeoutableProvider(p, option.RequestTimeout)
	p = NewRetriableProvider(p, option.RetryCount, option.RetryInterval)
	return p
}
