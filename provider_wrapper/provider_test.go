package providers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	for i := 0; i < 3; i++ {
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

func TestTimeoutProvider(t *testing.T) {
	// 创建一个模拟的慢速服务器
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, `{"jsonrpc":"2.0","id":1,"result":"ok"}`)
	}))
	defer slowServer.Close()

	// 创建一个超时时间为1秒的Provider
	p, err := NewProviderWithOption(slowServer.URL, Option{RequestTimeout: 1 * time.Second})
	assert.NilError(t, err)

	var result interface{}
	ctx := context.Background()

	// 调用应该因超时而失败
	err = p.CallContext(ctx, &result, "test_method")
	assert.Assert(t, err != nil, "预期出现超时错误")
	assert.ErrorContains(t, err, "timeout", "错误应该是超时")

	// 创建一个超时时间为3秒的Provider
	p, err = NewProviderWithOption(slowServer.URL, Option{RequestTimeout: 3 * time.Second})
	assert.NilError(t, err)

	// 调用应该成功
	err = p.CallContext(ctx, &result, "test_method")
	assert.NilError(t, err, "调用应该成功")
	assert.Equal(t, result, "ok", "结果应该是'ok'")
}
