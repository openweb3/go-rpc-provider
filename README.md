go-rpc-provider
===========
go-rpc-provider is enhanced on the basis of the package github.com/ethereum/go-ethereum/rpc

Features
-----------
-   Replace HTTP by fasthttp to improve performance of concurrency HTTP request
-   Refactor function `func (err *jsonError) Error() ` to return 'data' of JSON RPC response
-   Set HTTP request header 'Connection' to 'Keep-Alive' to reuse TCP connection for avoiding the amount of TIME_WAIT on the RPC server, and the RPC server also need to set 'Connection' to 'Keep-Alive'.
-   Add remote client address to WebSocket context for tracing.
-   Add client pointer to context when new RPC connection for tracing.
-   Add HTTP request in RPC context in function `ServeHTTP` for tracing.
-   Support variadic arguments for RPC service
-   Add provider wrapper for convenient extension
-   Support MiddlewarableProvider to extend provider features, currently support 
    -   `NewRetriableProvider` to create MiddlewarableProvider instance with auto-retry when failing
    -   `NewTimeoutableProvider` to create MiddlewarableProvider instance with timeout when CallContext/BatchCallContext
    -   `NewLoggerProvider` to create MiddlewarableProvider instance for logging request/response when CallContext/BatchCallContext
    -   `NewProviderWithOption` to create MiddlewarableProvider instance includes all features of `NewRetriableProvider`, `NewTimeoutableProvider` and `NewLoggerProvider`


Usage
-----------
rpc.Client implements a provider interface, so you could use it as a provider

create a simple RPC client
```golang
	rpc.Dial("http://localhost:8545")
```
create a simple RPC client with context for canceling or timeout the initial connection establishment
```golang
	rpc.DialContext("http://localhost:8545")
```

For feature extension, we apply MiddlewarableProvider for hook CallContext/BatchCallContext/Subscribe, such as log rpc request and response or cache environment variable in the context.

you can create MiddlewarableProvider by NewMiddlewarableProvider and pass the provider created below as the parameter

```golang
	p, e := rpc.DialContext(context.Background(), "http://localhost:8545")
	if e != nil {
		t.Fatal(e)
	}
	mp := NewMiddlewarableProvider(p)
	mp.HookCallContext(callContextLogMiddleware)
	mp.HookCallContext(otherMiddleware)
```

the callContextLogMiddleware is like
```golang
func callContextLogMiddleware(f providers.CallContextFunc) providers.CallContextFunc {
	return func(ctx context.Context, resultPtr interface{}, method string, args ...interface{}) error {
		fmt.Printf("request %v %v\n", method, args)
		err := f(ctx, resultPtr, method, args...)
		j, _ := json.Marshal(resultPtr)
		fmt.Printf("response %s\n", j)
		return err
	}
}
```

The simplest way to create middlewarable provider is by `providers.NewBaseProvider`.  It will create a MiddlewareProvider and it could set the max connection number for client with the server.

For example, created by `providers`.NewBaseProvider` and use the logging middleware created below to hook by HookCallContext
```golang
	p, e := warpper.NewBaseProvider(context.Background(), "http://localhost:8545")
	if e != nil {
		return e
	}
	p.HookCallContext(callLogMiddleware)
	NewClientWithProvider(p)
```

However, we also apply several functions to create kinds of instances of MiddlewarableProvider. 
The functions are `providers.NewTimeoutableProvider`, `providers.NewRetriableProvider`. 

And the simplest way to create NewMiddlewarableProvider with retry and timeout features is to use `providers.NewProviderWithOption`
<<<<<<< Updated upstream
=======


Server
----------

The `rpc` package is extended based on go-ethereum. And we extend RPC server features for hook on handle requests on both "batch call msg" and "call msg".

### Usage

Both `rpc.HookHandleMsg` and `rpc.HookHandleBatch` are globally effective, it will effective to both HTTP and Websocket when hook once.

#### Note:
- HookHandleCallMsg works when "call msg" and "batch call msg", for example, batch requests `[ jsonrpc_a, jsonrpc_b ]`, it will trigger function nested by `HookHandleBatch` once and function netsed by `HookHandleCallMsg` twice
- HookHandleBatch works only when "call msg"

```golang
	rpc.HookHandleCallMsg(func(next rpc.HandleCallMsgFunc) rpc.HandleCallMsgFunc {
		return func(ctx context.Context, msg *rpc.JsonRpcMessage) *rpc.JsonRpcMessage {
			fmt.Printf("request call msg %v\n", utils.PrettyJSON(msg))
			fmt.Printf("callmsg -- request context key of %v value %v\n", "test-k", ctx.Value("test-k"))
			result := next(ctx, msg)
			fmt.Printf("response call msg %v\n", utils.PrettyJSON(result))
			return result
		}
	})
	rpc.HookHandleBatch(func(next rpc.HandleBatchFunc) rpc.HandleBatchFunc {
		return func(ctx context.Context, msgs []*rpc.JsonRpcMessage) []*rpc.JsonRpcMessage {
			fmt.Printf("request batch %v\n", utils.PrettyJSON(msgs))
			fmt.Printf("batch -- request context key of %v value %v\n", "test-k", ctx.Value("test-k"))
			results := next(ctx, msgs)
			fmt.Printf("response batch %v\n", utils.PrettyJSON(results))
			return results
		}
	})
```
>>>>>>> Stashed changes
