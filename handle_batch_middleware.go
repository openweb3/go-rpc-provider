package rpc

import (
	"context"
)

type HandleBatchFunc func(ctx context.Context, msgs []*jsonrpcMessage) []*JsonRpcMessage
type HandleBatchMiddleware func(next HandleBatchFunc) HandleBatchFunc

var (
	handleBatchFuncMiddlewares []HandleBatchMiddleware
)

func HookHandleBatch(middleware HandleBatchMiddleware) {
	handleBatchFuncMiddlewares = append(handleBatchFuncMiddlewares, middleware)
}

func (h *handler) getHandleBatchNestedware() HandleBatchFunc {
	if h.handleBatchNestedWare == nil || h.handleBatchMiddlewareLen != len(handleBatchFuncMiddlewares) {
		h.handleBatchMiddlewareLen = len(handleBatchFuncMiddlewares)
		nestedWare := func(ctx context.Context, msgs []*jsonrpcMessage) []*JsonRpcMessage {
			c := h.handleBatchCore(ctx, msgs)
			if c == nil {
				return nil
			}
			result := <-c
			return result
		}
		for i := len(handleBatchFuncMiddlewares) - 1; i >= 0; i-- {
			nestedWare = handleBatchFuncMiddlewares[i](nestedWare)
		}
		h.handleBatchNestedWare = nestedWare
	}
	return h.handleBatchNestedWare
}
