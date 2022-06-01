package rpc

import (
	"context"
	"fmt"
)

type HandleBatchFunc func(ctx context.Context, msgs []*jsonrpcMessage) []*JsonRpcMessage
type HandleBatchMiddleware func(next HandleBatchFunc) HandleBatchFunc

var (
	handleBatchFuncMiddlewares []HandleBatchMiddleware
)

func HookHandleBatch(middleware HandleBatchMiddleware) {
	// fmt.Printf("**HookHandleBatch**")
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
			fmt.Printf("receive batch result from channel %v\n", result)
			return result
		}
		// fmt.Printf("len(handleBatchFuncMiddlewares) %v\n", len(handleBatchFuncMiddlewares))
		for i := len(handleBatchFuncMiddlewares) - 1; i >= 0; i-- {
			nestedWare = handleBatchFuncMiddlewares[i](nestedWare)
		}
		h.handleBatchNestedWare = nestedWare
	}
	return h.handleBatchNestedWare
}

// func handleBatchMid(next HandleBatchFunc) HandleBatchFunc {
// 	return func(msgs []*jsonrpcMessage) {
// 		// do sth pre
// 		next(msgs)
// 		// do sth post
// 	}
// }
