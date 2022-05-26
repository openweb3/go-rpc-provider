package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/openweb3/go-rpc-provider"
	"gotest.tools/assert"
)

var executeStack []byte = make([]byte, 0)

func TestHookCallContext(t *testing.T) {
	executeStack = make([]byte, 0)
	p, e := rpc.DialContext(context.Background(), "http://localhost:8545")
	if e != nil {
		t.Fatal(e)
	}
	mp := NewMiddlewarableProvider(p)

	mp.HookCallContext(callContextMiddle1)
	mp.HookCallContext(callContextMiddle2)

	b := new(hexutil.Big)
	mp.CallContext(context.Background(), b, "eth_blockNumber")
	assert.DeepEqual(t, executeStack, []byte{1, 2, 3, 4})
}

func callContextMiddle1(f CallContextFunc) CallContextFunc {
	return func(ctx context.Context, resultPtr interface{}, method string, args ...interface{}) error {
		executeStack = append(executeStack, 1)
		fmt.Printf("call %v %v\n", method, args)
		ctx = context.WithValue(ctx, "foo", "bar")
		err := f(ctx, resultPtr, method, args...)
		j, _ := json.Marshal(resultPtr)
		fmt.Printf("response %s\n", j)
		executeStack = append(executeStack, 4)
		return err
	}
}

func callContextMiddle2(f CallContextFunc) CallContextFunc {
	return func(ctx context.Context, resultPtr interface{}, method string, args ...interface{}) error {
		executeStack = append(executeStack, 2)
		fmt.Printf("call %v %v with context key foo value %v\n", method, args, ctx.Value("foo"))
		err := f(ctx, resultPtr, method, args...)
		executeStack = append(executeStack, 3)
		return err
	}
}
