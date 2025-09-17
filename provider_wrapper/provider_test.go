package providers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
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
	p, err := NewProviderWithOption("http://localhost:12345", *new(Option).WithCircuitBreaker(DefaultCircuitBreakerOption{
		MaxFail:        1,
		FailTimeWindow: time.Second * 2,
		OpenColdTime:   time.Second,
	}))
	assert.NilError(t, err)

	var result interface{}

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			for j := 0; j < 6; j++ {
				err = p.CallContext(context.Background(), &result, "xx")
				fmt.Printf("%v %v\n", time.Now().Format(time.RFC3339), err)
				time.Sleep(time.Millisecond * 500)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
func TestTimeoutProvider(t *testing.T) {
	// Create a mock slow server
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, `{"jsonrpc":"2.0","id":1,"result":"ok"}`)
	}))
	defer slowServer.Close()

	// Create a Provider with 1 second timeout
	p, err := NewProviderWithOption(slowServer.URL, Option{RequestTimeout: 1 * time.Second})
	assert.NilError(t, err)

	var result interface{}
	ctx := context.Background()

	// Call should fail due to timeout
	err = p.CallContext(ctx, &result, "test_method")
	assert.Assert(t, err != nil, "Expected timeout error")
	assert.ErrorContains(t, err, "timeout", "Error should be timeout")

	// Create a Provider with 3 second timeout
	p, err = NewProviderWithOption(slowServer.URL, Option{RequestTimeout: 3 * time.Second})
	assert.NilError(t, err)

	// Call should succeed
	err = p.CallContext(ctx, &result, "test_method")
	assert.NilError(t, err, "Call should succeed")
	assert.Equal(t, result, "ok", "Result should be 'ok'")
}
