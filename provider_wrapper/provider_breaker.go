package providers

import (
	"context"

	"github.com/openweb3/go-rpc-provider"
	"github.com/openweb3/go-rpc-provider/interfaces"
)

func NewBreakerProvider(inner interfaces.Provider, breaker ICircuitBreaker) *MiddlewarableProvider {
	m := NewMiddlewarableProvider(inner)

	r := &BreakerMiddleware{breaker}
	m.HookCallContext(r.callContextMiddleware)
	m.HookBatchCallContext(r.batchCallContextMiddleware)
	return m
}

type BreakerMiddleware struct {
	breaker ICircuitBreaker
}

func (r *BreakerMiddleware) callContextMiddleware(call CallContextFunc) CallContextFunc {
	return func(ctx context.Context, resultPtr interface{}, method string, args ...interface{}) error {
		handler := func() error {
			return call(ctx, resultPtr, method, args...)
		}
		return r.breaker.Do(handler)
	}
}

func (r *BreakerMiddleware) batchCallContextMiddleware(call BatchCallContextFunc) BatchCallContextFunc {
	return func(ctx context.Context, b []rpc.BatchElem) error {
		handler := func() error {
			return call(ctx, b)
		}
		return r.breaker.Do(handler)
	}
}
