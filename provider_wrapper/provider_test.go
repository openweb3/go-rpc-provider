package providers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mcuadros/go-defaults"
	"gotest.tools/assert"
)

func TestConfigurationDefault(t *testing.T) {
	c := Option{}
	defaults.SetDefaults(&c)
	assert.Equal(t, c.RetryCount, 0)
	assert.Equal(t, c.RetryInterval, 1*time.Second)
	assert.Equal(t, c.RequestTimeout, 30*time.Second)

	c = Option{RetryCount: 10}
	defaults.SetDefaults(&c)
	assert.Equal(t, c.RetryCount, 10)
	assert.Equal(t, c.RetryInterval, 1*time.Second)
	assert.Equal(t, c.RequestTimeout, 30*time.Second)
}

func TestProviderShouldCircuitBreak(t *testing.T) {
	p, err := NewProviderWithOption("http://localhost:1234", *new(Option).WithCircuitBreaker(DefaultCircuitBreakerOption{}))
	assert.NilError(t, err)

	var result interface{}
	for i := 0; i < 1; i++ {
		go func() {
			for {
				err = p.CallContext(context.Background(), &result, "xx")
				fmt.Printf("%v %v\n", time.Now().Format(time.RFC3339), err)
				time.Sleep(time.Millisecond * 500)
			}
		}()
	}
	select {}
}
