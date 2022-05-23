package providers

import (
	"context"

	"github.com/openweb3/go-rpc-provider"
	"github.com/openweb3/go-rpc-provider/interfaces"
)

type MiddlewarableProvider struct {
	Inner interfaces.Provider
	// callMiddles             []CallMiddleware
	callContextMiddles []CallContextMiddleware
	// batchCallMiddles        []BatchCallMiddleware
	batchCallContextMiddles []BatchCallContextMiddleware
	subscribeMiddles        []SubscribeMiddleware

	// callNestedWare             CallFunc
	callContextNestedWare CallContextFunc
	// batchcallNestedWare        BatchCallFunc
	batchcallContextNestedWare BatchCallContextFunc
	subscribeNestedWare        SubscribeFunc
}

func NewMiddlewarableProvider(p interfaces.Provider) *MiddlewarableProvider {
	m := &MiddlewarableProvider{Inner: p,
		// callNestedWare:             p.Call,
		callContextNestedWare: p.CallContext,
		// batchcallNestedWare:        p.BatchCall,
		batchcallContextNestedWare: p.BatchCallContext,
		subscribeNestedWare:        p.Subscribe,
	}
	return m
}

// type CallFunc func(resultPtr interface{}, method string, args ...interface{}) error
type CallContextFunc func(ctx context.Context, result interface{}, method string, args ...interface{}) error

// type BatchCallFunc func(b []rpc.BatchElem) error
type BatchCallContextFunc func(ctx context.Context, b []rpc.BatchElem) error
type SubscribeFunc func(ctx context.Context, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error)

// type SubscribeFunc func(ctx context.Context, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error)

// type CallMiddleware func(CallFunc) CallFunc
type CallContextMiddleware func(CallContextFunc) CallContextFunc

// type BatchCallMiddleware func(BatchCallFunc) BatchCallFunc
type BatchCallContextMiddleware func(BatchCallContextFunc) BatchCallContextFunc
type SubscribeMiddleware func(SubscribeFunc) SubscribeFunc

// type SubscribeMiddleware func(SubscribeFunc) SubscribeFunc

// callMiddles: A -> B -> C, execute order is A -> B -> C
// func (mp *MiddlewarableProvider) HookCall(cm CallMiddleware) {
// 	mp.callMiddles = append(mp.callMiddles, cm)
// 	mp.callNestedWare = mp.Inner.Call
// 	for i := len(mp.callMiddles) - 1; i >= 0; i-- {
// 		mp.callNestedWare = mp.callMiddles[i](mp.callNestedWare)
// 	}
// }

func (mp *MiddlewarableProvider) HookCallContext(cm CallContextMiddleware) {
	mp.callContextMiddles = append(mp.callContextMiddles, cm)
	mp.callContextNestedWare = mp.Inner.CallContext
	for i := len(mp.callContextMiddles) - 1; i >= 0; i-- {
		mp.callContextNestedWare = mp.callContextMiddles[i](mp.callContextNestedWare)
	}
}

// func (mp *MiddlewarableProvider) HookBatchCall(cm BatchCallMiddleware) {
// 	mp.batchCallMiddles = append(mp.batchCallMiddles, cm)
// 	mp.batchcallNestedWare = mp.Inner.BatchCall
// 	for i := len(mp.batchCallMiddles) - 1; i >= 0; i-- {
// 		mp.batchcallNestedWare = mp.batchCallMiddles[i](mp.batchcallNestedWare)
// 	}

// }

func (mp *MiddlewarableProvider) HookBatchCallContext(cm BatchCallContextMiddleware) {
	mp.batchCallContextMiddles = append(mp.batchCallContextMiddles, cm)
	mp.batchcallContextNestedWare = mp.Inner.BatchCallContext
	for i := len(mp.batchCallContextMiddles) - 1; i >= 0; i-- {
		mp.batchcallContextNestedWare = mp.batchCallContextMiddles[i](mp.batchcallContextNestedWare)
	}
}

func (mp *MiddlewarableProvider) HookSubscribe(cm SubscribeMiddleware) {
	mp.subscribeMiddles = append(mp.subscribeMiddles, cm)
	mp.subscribeNestedWare = mp.Inner.Subscribe
	for i := len(mp.subscribeMiddles) - 1; i >= 0; i-- {
		mp.subscribeNestedWare = mp.subscribeMiddles[i](mp.subscribeNestedWare)
	}
}

// func (mp *MiddlewarableProvider) Call(resultPtr interface{}, method string, args ...interface{}) error {
// 	return mp.callNestedWare(resultPtr, method, args...)
// }

func (mp *MiddlewarableProvider) CallContext(ctx context.Context, resultPtr interface{}, method string, args ...interface{}) error {
	return mp.callContextNestedWare(ctx, resultPtr, method, args...)
}

// func (mp *MiddlewarableProvider) BatchCall(b []rpc.BatchElem) error {
// 	return mp.batchcallNestedWare(b)
// }

func (mp *MiddlewarableProvider) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return mp.batchcallContextNestedWare(ctx, b)
}

func (mp *MiddlewarableProvider) Subscribe(ctx context.Context, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error) {
	return mp.subscribeNestedWare(ctx, namespace, channel, args...)
}

func (mp *MiddlewarableProvider) Close() {
	mp.Inner.Close()
}
