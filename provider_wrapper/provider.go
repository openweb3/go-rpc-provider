package providers

import (
	"context"
	"time"

	"github.com/mcuadros/go-defaults"
	"github.com/openweb3/go-rpc-provider/interfaces"
)

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

func NewProviderWithOption(rawurl string, option *Option) (interfaces.Provider, error) {
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
func wrapProvider(p interfaces.Provider, option *Option) interfaces.Provider {
	if option == nil {
		return p
	}

	p = NewTimeoutableProvider(p, option.RequestTimeout)
	p = NewRetriableProvider(p, option.RetryCount, option.RetryInterval)
	return p
}
