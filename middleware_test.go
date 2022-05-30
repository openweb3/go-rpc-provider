package rpc

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestHandleMsgMiddleware(t *testing.T) {
	// set new handler hook
	fmt.Println("start")
	HookContextOnNewHandler(onNewHandlerMiddleware1)
	HookContextOnNewHandler(onNewHandlerMiddleware2)

	// new http server
	server := newTestServer()
	defer server.Stop()
	client := DialInProc(server)
	defer client.Close()

	var resp echoResult
	if err := client.Call(&resp, "test_echo", "hello", 10, &echoArgs{"world"}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(resp, echoResult{"hello", 10, &echoArgs{"world"}}) {
		t.Errorf("incorrect result %#v", resp)
	}
}

func onNewHandlerMiddleware1(next SetContextFunc) SetContextFunc {
	return func(ctx context.Context) context.Context {
		ctx = context.WithValue(ctx, "foo", "bar")
		next(ctx)
		return ctx
	}
}

func onNewHandlerMiddleware2(next SetContextFunc) SetContextFunc {
	return func(ctx context.Context) context.Context {
		fmt.Printf("ctx foo value: %v\n", ctx.Value("foo"))
		next(ctx)
		return ctx
	}
}
