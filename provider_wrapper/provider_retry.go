package providers

import (
	"context"
	"time"

	"github.com/openweb3/go-rpc-provider"
	"github.com/openweb3/go-rpc-provider/interfaces"
	"github.com/openweb3/go-rpc-provider/utils"
	"github.com/pkg/errors"
)

func NewRetriableProvider(inner interfaces.Provider, maxRetry int, interval time.Duration) *MiddlewarableProvider {
	m := NewMiddlewarableProvider(inner)

	r := &RetriableMiddleware{maxRetry, interval}
	m.HookCall(r.callMiddleware)
	m.HookCallContext(r.callContextMiddleware)
	m.HookBatchCall(r.batchCallMiddleware)
	m.HookBatchCallContext(r.batchCallContextMiddleware)
	return m
}

type RetriableMiddleware struct {
	maxRetry int
	interval time.Duration
}

func (r *RetriableMiddleware) callMiddleware(call CallFunc) CallFunc {
	return func(resultPtr interface{}, method string, args ...interface{}) error {
		handler := func() error {
			return call(resultPtr, method, args...)
		}
		return retry(r.maxRetry, r.interval, handler)
	}
}

func (r *RetriableMiddleware) callContextMiddleware(call CallContextFunc) CallContextFunc {
	return func(ctx context.Context, resultPtr interface{}, method string, args ...interface{}) error {
		handler := func() error {
			return call(ctx, resultPtr, method, args...)
		}
		return retry(r.maxRetry, r.interval, handler)
	}
}

func (r *RetriableMiddleware) batchCallMiddleware(call BatchCallFunc) BatchCallFunc {
	return func(b []rpc.BatchElem) error {
		handler := func() error {
			return call(b)
		}
		return retry(r.maxRetry, r.interval, handler)
	}
}

func (r *RetriableMiddleware) batchCallContextMiddleware(call BatchCallContextFunc) BatchCallContextFunc {
	return func(ctx context.Context, b []rpc.BatchElem) error {
		handler := func() error {
			return call(ctx, b)
		}
		return retry(r.maxRetry, r.interval, handler)
	}
}

func retry(maxRetry int, interval time.Duration, handler func() error) error {
	remain := maxRetry
	for {
		err := handler()
		if err == nil {
			return nil
		}

		if utils.IsRPCJSONError(err) {
			return err
		}

		if remain <= 0 {
			return errors.Wrapf(err, "failed after %v retries", maxRetry)
		}

		if interval > 0 {
			time.Sleep(interval)
		}

		remain--
	}
}
