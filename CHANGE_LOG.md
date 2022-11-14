# Change Log

## v0.3.1
- Support block tag: finalize and safe 

## v0.2.1
- Support logger middleware
- Set default retry count to 0
- Set default request timeout to 30s

## v0.2.0
- Remove Call/BatchCall and HookCall/HookBatchall for avoiding user confuse
- Unify all created provider to MiddlewareProvider in providers package
- Support NewBaseProvider/NewTimeoutableProvider/NewRetriableProvider/NewProviderWithOption and all of them return a MiddlewarableProvider