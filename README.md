go-rpc-provider
===========
go-rpc-provider is enhanced on the basis of the package github.com/ethereum/go-ethereum/rpc

Features
-----------
-   Replace http by fasthttp to imporve performance of concurrency http request
-   Refactor function `func (err *jsonError) Error() ` to return 'data' of Json RPC response
-   Set http request header 'Connection' to 'Keep-Alive' to reuse tcp connection for avoiding amount of TIME_WAIT on rpc server, and the rpc server also need set 'Connection' to 'Keep-Alive'.
-   Add remote client addrto websocket context for tracing.
-   Add client pointer to context when new rpc connection for tracing.
-   Add http request in RPC context in fucntion `ServeHTTP` for tracing.
-   Support variadic arguments for rpc service