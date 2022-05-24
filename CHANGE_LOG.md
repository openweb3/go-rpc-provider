# Change Log

## v1.1.0
- Remove Call/BatchCall and HookCall/HookBatchall for avoiding user confuse
- Unify all created provider to MiddlewareProvider in providers package
- Support NewBaseProvider/NewTimeoutableProvider/NewRetriableProvider/NewProviderWithOption and all of them return a MiddlewarableProvider