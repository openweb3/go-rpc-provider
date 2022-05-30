package rpc

import "context"

type JsonRpcMessage = jsonrpcMessage

// Handle New Handler
var (
	onNewHandlerMiddlewares []NewHandlerMiddleware
	onNewHandlerNestedWare  SetContextFunc = defaultSetContext
)

type SetContextFunc func(connCtx context.Context) context.Context
type NewHandlerMiddleware func(next SetContextFunc) SetContextFunc

func HookContextOnNewHandler(next NewHandlerMiddleware) {
	onNewHandlerMiddlewares = append(onNewHandlerMiddlewares, next)
	onNewHandlerNestedWare = defaultSetContext
	for i := len(onNewHandlerMiddlewares) - 1; i >= 0; i-- {
		onNewHandlerNestedWare = onNewHandlerMiddlewares[i](onNewHandlerNestedWare)
	}
}

func defaultSetContext(connCtx context.Context) context.Context {
	return connCtx
}

// Handle Msg Middleware

type HandleMsgFunc func(msg *JsonRpcMessage) *JsonRpcMessage
type HandleMsgMiddleware func(next HandleMsgFunc) HandleMsgFunc

var (
	handleMsgFuncMiddlewares []HandleMsgMiddleware
)

func HookHandleMsg(middleware HandleMsgMiddleware) {
	handleMsgFuncMiddlewares = append(handleMsgFuncMiddlewares, middleware)
}

func (h *handler) getHandleMsgNestedware() HandleMsgFunc {
	if h.handleMsgNestedware == nil || h.handleMsgMiddlewareLen != len(handleMsgFuncMiddlewares) {
		h.handleMsgMiddlewareLen = len(handleMsgFuncMiddlewares)
		nestedWare := func(msg *jsonrpcMessage) *JsonRpcMessage {
			c := h.handleMsgCore(msg)
			if c == nil {
				return nil
			}
			result := <-c
			return result
		}
		for i := len(handleMsgFuncMiddlewares) - 1; i >= 0; i-- {
			nestedWare = handleMsgFuncMiddlewares[i](nestedWare)
		}
		h.handleMsgNestedware = nestedWare
	}
	return h.handleMsgNestedware
}
