package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/openweb3/go-rpc-provider"
	"github.com/openweb3/go-rpc-provider/interfaces"
	"github.com/openweb3/go-rpc-provider/utils"
)

type LogMiddleware struct {
	Writer io.Writer
}

func NewLoggerProvider(p interfaces.Provider, w io.Writer) *MiddlewarableProvider {
	mp := NewMiddlewarableProvider(p)
	logMiddleware := &LogMiddleware{Writer: w}
	mp.HookCallContext(logMiddleware.callContextMiddleware)
	mp.HookBatchCallContext(logMiddleware.batchCallContextMiddleware)
	return mp
}

func (m *LogMiddleware) callContextMiddleware(call CallContextFunc) CallContextFunc {
	return func(ctx context.Context, resultPtr interface{}, method string, args ...interface{}) error {
		argsStr := fmt.Sprintf("%+v", args)
		argsJson, err := json.Marshal(args)
		if err == nil {
			argsStr = string(argsJson)
		}

		start := time.Now()
		err = call(ctx, resultPtr, method, args...)
		duration := time.Since(start)

		if m.Writer == os.Stdout {
			if err == nil {
				fmt.Printf("%v Method %v, Params %v, Result %v, Use %v\n",
					color.GreenString("[Call RPC Done]"),
					color.YellowString(method),
					color.CyanString(argsStr),
					color.CyanString(utils.PrettyJSON(resultPtr)),
					color.CyanString(duration.String()))
				return nil
			}

			color.Red("%v Method %v, Params %v, Error %v, Use %v\n",
				color.RedString("[Call RPC Fail]"),
				color.YellowString(method),
				color.CyanString(string(argsJson)),
				color.RedString(fmt.Sprintf("%+v", err)),
				color.CyanString(duration.String()))
			return err
		}

		if err == nil {
			l := fmt.Sprintf("%v Method %v, Params %v, Result %v, Use %v\n",
				"[Call RPC Done]", method, argsStr, utils.PrettyJSON(resultPtr), duration.String())
			m.Writer.Write([]byte(l))
			return nil
		}

		el := fmt.Sprintf("%v Method %v, Params %v, Error %v, Use %v\n",
			"[Call RPC Fail]", method, string(argsJson), fmt.Sprintf("%+v", err), duration.String())
		m.Writer.Write([]byte(el))
		return err
	}
}

func (m *LogMiddleware) batchCallContextMiddleware(batchCall BatchCallContextFunc) BatchCallContextFunc {
	return func(ctx context.Context, b []rpc.BatchElem) error {
		start := time.Now()

		err := batchCall(ctx, b)

		duration := time.Since(start)

		if m.Writer == os.Stdout {
			if err == nil {
				fmt.Printf("%v BatchElems %v, Use %v\n",
					color.GreenString("[Batch Call RPC Done]"),
					color.CyanString(utils.PrettyJSON(b)),
					color.CyanString(duration.String()))
				return nil
			}
			fmt.Printf("%v BatchElems %v, Error: %v, Use %v\n",
				color.RedString("[Batch Call RPC Fail]"),
				color.CyanString(utils.PrettyJSON(b)),
				color.RedString(fmt.Sprintf("%+v", err)),
				duration)
			return err
		}

		if err == nil {
			l := fmt.Sprintf("%v BatchElems %v, Use %v\n",
				"[Batch Call RPC Done]",
				utils.PrettyJSON(b),
				duration.String())
			m.Writer.Write([]byte(l))
			return nil
		}

		el := fmt.Sprintf("%v BatchElems %v, Error: %v, Use %v\n",
			"[Batch Call RPC Fail]",
			utils.PrettyJSON(b),
			fmt.Sprintf("%+v", err),
			duration)
		m.Writer.Write([]byte(el))
		return err
	}
}
