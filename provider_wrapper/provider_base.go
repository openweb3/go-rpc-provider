package providers

import (
	"context"
	"net/url"

	"github.com/openweb3/go-rpc-provider"

	"github.com/valyala/fasthttp"
)

// NewBaseProvider returns a new BaseProvider.
// maxConnsPerHost is the maximum number of concurrent connections, only works for http(s) protocal
func NewBaseProvider(ctx context.Context, nodeUrl string, maxConnectionNum ...int) (*MiddlewarableProvider, error) {
	if len(maxConnectionNum) > 0 && maxConnectionNum[0] > 0 {
		if u, err := url.Parse(nodeUrl); err == nil {
			if u.Scheme == "http" || u.Scheme == "https" {
				fasthttpClient := new(fasthttp.Client)
				fasthttpClient.MaxConnsPerHost = maxConnectionNum[0]
				p, err := rpc.DialHTTPWithClient(nodeUrl, fasthttpClient)
				if err != nil {
					return nil, err
				}
				return NewMiddlewarableProvider(p), nil
			}
		}
	}

	p, err := rpc.DialContext(ctx, nodeUrl)
	if err != nil {
		return nil, err
	}
	return NewMiddlewarableProvider(p), nil
}

// MustNewBaseProvider returns a new BaseProvider. Panic if error.
func MustNewBaseProvider(ctx context.Context, nodeUrl string, maxConnectionNum ...int) *MiddlewarableProvider {
	p, err := NewBaseProvider(ctx, nodeUrl, maxConnectionNum...)
	if err != nil {
		panic(err)
	}
	return p
}
