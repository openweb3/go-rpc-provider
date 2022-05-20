package providers

import (
	"context"
	"time"

	rpc "github.com/openweb3/go-rpc-provider"
	"github.com/openweb3/go-rpc-provider/interfaces"
)

// TimeoutableProvider overwrite Call by CallContext with timeout context, make it to internal package to prevent external package to use it.
type TimeoutableProvider struct {
	MiddlewarableProvider
}

func NewTimeoutableProvider(inner interfaces.Provider, timeout time.Duration) *MiddlewarableProvider {
	m := NewMiddlewarableProvider(inner)
	timeoutMiddle := TimeoutMiddleware{Timeout: timeout}

	m.HookCallContext(timeoutMiddle.CallContext)
	m.HookBatchCallContext(timeoutMiddle.BatchCallContext)
	m.HookSubscribe(timeoutMiddle.Subscribe)

	return m
}

func (p *TimeoutableProvider) Call(resultPtr interface{}, method string, args ...interface{}) error {
	return p.CallContext(context.Background(), resultPtr, method, args...)
}

func (p *TimeoutableProvider) BatchCall(b []rpc.BatchElem) error {
	return p.BatchCallContext(context.Background(), b)
}

type TimeoutMiddleware struct {
	Timeout time.Duration
}

func (t TimeoutMiddleware) CallContext(call CallContextFunc) CallContextFunc {
	return func(ctx context.Context, resultPtr interface{}, method string, args ...interface{}) error {
		ctx, f := t.setContext(ctx)
		defer f()
		return call(ctx, resultPtr, method, args...)
	}
}

func (t TimeoutMiddleware) BatchCallContext(call BatchCallContextFunc) BatchCallContextFunc {
	return func(ctx context.Context, b []rpc.BatchElem) error {
		ctx, f := t.setContext(ctx)
		defer f()
		return call(ctx, b)
	}
}

func (t TimeoutMiddleware) Subscribe(sub SubscribeFunc) SubscribeFunc {
	return func(ctx context.Context, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error) {
		ctx, f := t.setContext(ctx)
		defer f()
		return sub(ctx, namespace, channel, args...)
	}
}

func (t *TimeoutMiddleware) setContext(ctx context.Context) (context.Context, context.CancelFunc) {
	_, ok := ctx.Deadline()
	if !ok {
		ctx, f := context.WithTimeout(ctx, t.Timeout)
		return ctx, f
	}
	return ctx, func() {}
}
