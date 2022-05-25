package utils

import (
	"bytes"
	"encoding/json"
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

// PrettyJSON json marshal value and pretty with indent
func PrettyJSON(value interface{}) string {
	j, e := json.Marshal(value)
	if e != nil {
		panic(e)
	}
	var str bytes.Buffer
	_ = json.Indent(&str, j, "", "    ")
	return str.String()
}
