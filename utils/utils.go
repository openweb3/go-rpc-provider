package utils

import (
	"reflect"

	"github.com/pkg/errors"
)

// IsRPCJSONError returns true if err is rpc error
func IsRPCJSONError(err error) bool {
	if err == nil {
		return false
	}

	t := reflect.TypeOf(errors.Cause(err)).String()
	return t == "*rpc.jsonError" || t == "rpc.jsonError" || t == "utils.RpcError" || t == "*utils.RpcError"
}
