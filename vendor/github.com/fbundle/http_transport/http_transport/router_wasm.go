//go:build js && wasm

package http_transport

import (
	"strings"
	"syscall/js"
)

// WasmRouterGroup mirrors gin.RouterGroup but registers handlers as JS global
// functions via js.Global().Set.
//
// Path segments are joined with "$": Group("sudoku").POST("api/new") →
// window.POST$sudoku$api$new(jsonIn) → jsonOut.
type WasmRouterGroup struct{ prefix string }

// New returns a root WasmRouterGroup (no prefix).
func NewWASM() *WasmRouterGroup { return &WasmRouterGroup{} }

var New = NewWASM

// Group returns a child WasmRouterGroup with the given path segment appended.
func (r *WasmRouterGroup) Group(path string) *WasmRouterGroup {
	return &WasmRouterGroup{prefix: join(r.prefix, path)}
}

func httpTransportJsName(method, path string) string {
	return method + "$" + strings.ReplaceAll(path, "/", "$")
}

// POST registers fn as a JS global function reachable from JavaScript.
func (r *WasmRouterGroup) POST(relativePath string, fn HandlerFunc) {
	name := httpTransportJsName("POST", join(r.prefix, relativePath))
	js.Global().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
		var body []byte
		if len(args) > 0 {
			body = []byte(args[0].String())
		}
		_, resp := fn(body)
		return string(resp)
	}))
}
