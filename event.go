package rpc

import (
	"context"

	"github.com/valyala/fasthttp"
)

var beforeSendHttpHandlers []BeforeSendHttp

type BeforeSendHttp func(ctx context.Context, req *fasthttp.Request) error

func RegisterBeforeSendHttp(fn BeforeSendHttp) {
	if fn != nil {
		beforeSendHttpHandlers = append(beforeSendHttpHandlers, fn)
	}
}
