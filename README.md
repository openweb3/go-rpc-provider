go-rpc-provider
===========
go-rpc-provider is enhanced on the basis of the package github.com/ethereum/go-ethereum/rpc

Features
-----------
-   Replace http by fasthttp to imporve performance of concurrency http request
-   Refactor function `func (err *jsonError) Error() ` to return 'data' of Json RPC response
-   Set http request header 'Connection' to 'Keep-Alive' to reuse tcp connection for avoiding amount of TIME_WAIT on rpc server, and the rpc server also need set 'Connection' to 'Keep-Alive'.
-   Add remote client address to websocket context for tracing.
-   Add client pointer to context when new rpc connection for tracing.
-   Add http request in RPC context in fucntion `ServeHTTP` for tracing.
-   Support variadic arguments for rpc service
-   Add provider wraaper for convinent extension



Usage
-----------
rpc.Client implements provider interface, so you could use it as a provider

create simple rpc client
```golang
    rpc.Dial("http://localhost:8545")
```
create simple rpc client with context for cancel or timeout the initial connection establishment
```golang
	rpc.DialContext("http://localhost:8545")
```
or create base provider, it's a simple wrapper same as create client but it could set max connection number for client with server
```golang
	p, e := providers.NewBaseProvider(context.Background(), "http://localhost:8545", 1000)
	if e != nil {
		return e
	}
	NewClientWithProvider(p)
```
for custom pre/post call/batchcall behaviors, you can use `NewRetriableProvider` to create MiddlewareProvider for hooking on call/batchcall, such as log requests and so on
```golang
	p, e := warpper.NewBaseProvider(context.Background(), "http://localhost:8545")
	if e != nil {
		return e
	}
	mp := providers.NewMiddlewarableProvider(p)
	mp.HookCall(callLogMiddleware)
	NewClientWithProvider(p)
```
the callLogMiddleware is like
```golang
func callLogMiddleware(f providers.CallFunc) providers.CallFunc {
	return func(resultPtr interface{}, method string, args ...interface{}) error {
		fmt.Printf("request %v %v\n", method, args)
		err := f(resultPtr, method, args...)
		j, _ := json.Marshal(resultPtr)
		fmt.Printf("response %s\n", j)
		return err
	}
}
```
use `NewRetriableProvider` to create provider with retry and timeout feature, it will retry when call/batchcall failed and timeout in specific times
```golang
	p, e := providers.NewBaseProvider(context.Background(), "http://localhost:8545")
		if e != nil {
			return e
		}
	NewRetriableProvider(p,3,time.Second)
```