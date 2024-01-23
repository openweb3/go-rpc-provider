package providers

import (
	"context"

	"github.com/openweb3/go-rpc-provider"
	"github.com/openweb3/go-rpc-provider/interfaces"
)

func NewCircuitBreakerProvider(inner interfaces.Provider, breaker CircuitBreaker) *MiddlewarableProvider {
	m := NewMiddlewarableProvider(inner)

	b := &CircuitBreakerMiddleware{breaker}
	m.HookCallContext(b.callContextMiddleware)
	m.HookBatchCallContext(b.batchCallContextMiddleware)
	return m
}

type CircuitBreakerMiddleware struct {
	breaker CircuitBreaker
}

func (r *CircuitBreakerMiddleware) callContextMiddleware(call CallContextFunc) CallContextFunc {
	return func(ctx context.Context, resultPtr interface{}, method string, args ...interface{}) error {
		handler := func() error {
			return call(ctx, resultPtr, method, args...)
		}
		return r.breaker.Do(handler)
	}
}

func (r *CircuitBreakerMiddleware) batchCallContextMiddleware(call BatchCallContextFunc) BatchCallContextFunc {
	return func(ctx context.Context, b []rpc.BatchElem) error {
		handler := func() error {
			return call(ctx, b)
		}
		return r.breaker.Do(handler)
	}
}
