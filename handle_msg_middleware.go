package rpc

import "context"

type JsonRpcMessage = jsonrpcMessage

// Handle CallMsg Middleware
type HandleCallMsgFunc func(ctx context.Context, msg *JsonRpcMessage) *JsonRpcMessage
type HandleCallMsgMiddleware func(next HandleCallMsgFunc) HandleCallMsgFunc

var (
	handleCallMsgFuncMiddlewares []HandleCallMsgMiddleware
)

func HookHandleCallMsg(middleware HandleCallMsgMiddleware) {
	handleCallMsgFuncMiddlewares = append(handleCallMsgFuncMiddlewares, middleware)
}

func (h *handler) getHandleCallMsgNestedware(cp *callProc) HandleCallMsgFunc {
	nestedWare := func(ctx context.Context, msg *jsonrpcMessage) *JsonRpcMessage {
		cp.ctx = ctx
		return h.handleCallMsgCore(cp, msg)
	}
	for i := len(handleCallMsgFuncMiddlewares) - 1; i >= 0; i-- {
		nestedWare = handleCallMsgFuncMiddlewares[i](nestedWare)
	}
	return nestedWare
}
