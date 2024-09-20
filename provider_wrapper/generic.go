package providers

import (
	"context"

	"github.com/openweb3/go-rpc-provider"
	"github.com/openweb3/go-rpc-provider/interfaces"
)

// Call is a generic method to call RPC.
func Call[T any](provider interfaces.Provider, method string, args ...any) (result T, err error) {
	return CallContext[T](provider, context.Background(), method, args...)
}

// CallContext is a generic method to call RPC with context.
func CallContext[T any](provider interfaces.Provider, ctx context.Context, method string, args ...any) (result T, err error) {
	err = provider.CallContext(ctx, &result, method, args...)
	return
}

type Request struct {
	Method string
	Args   []any
}

type Response[T any] struct {
	Data  T
	Error error
}

// BatchCallContext is a generic method to call RPC with context in batch.
func BatchCallContext[T any](provider interfaces.Provider, ctx context.Context, requests ...Request) ([]Response[T], error) {
	batch := make([]rpc.BatchElem, 0, len(requests))
	responses := make([]Response[T], len(requests))

	for i, v := range requests {
		batch = append(batch, rpc.BatchElem{
			Method: v.Method,
			Args:   v.Args,
			Result: &responses[i].Data,
		})
	}

	if err := provider.BatchCallContext(ctx, batch); err != nil {
		return nil, err
	}

	for i, v := range batch {
		if v.Error != nil {
			responses[i].Error = v.Error
		}
	}

	return responses, nil
}

// BatchCallOneContext is generic method to call a single RPC with context in batch.
func BatchCallOneContext[T any](provider interfaces.Provider, ctx context.Context, method string, batchArgs ...[]any) ([]Response[T], error) {
	requests := make([]Request, 0, len(batchArgs))

	for _, v := range batchArgs {
		requests = append(requests, Request{
			Method: method,
			Args:   v,
		})
	}

	return BatchCallContext[T](provider, ctx, requests...)
}

// BatchCall is a generic method to call RPC in batch.
func BatchCall[T any](provider interfaces.Provider, requests ...Request) ([]Response[T], error) {
	return BatchCallContext[T](provider, context.Background(), requests...)
}

// BatchCallOne is a generic method to call a single RPC in batch.
func BatchCallOne[T any](provider interfaces.Provider, method string, batchArgs ...[]any) ([]Response[T], error) {
	return BatchCallOneContext[T](provider, context.Background(), method, batchArgs...)
}
