package providers

import "github.com/openweb3/go-rpc-provider/interfaces"

func IsMiddlewareableProvider(p interfaces.Provider) bool {
	_, ok := p.(*MiddlewarableProvider)
	return ok
}
