//go:build !js

package http_transport

import (
	"io"
	"net/http"
)

// GoRouterGroup implements Router using the standard library's net/http — no
// external dependencies required.
type GoRouterGroup struct {
	prefix string
	mux    *http.ServeMux
}

// NewGo wraps an existing http.ServeMux.
func NewGo(mux *http.ServeMux) *GoRouterGroup { return &GoRouterGroup{mux: mux} }

// Group returns a child GoRouterGroup with the given path segment appended.
func (r *GoRouterGroup) Group(path string) *GoRouterGroup {
	return &GoRouterGroup{prefix: join(r.prefix, path), mux: r.mux}
}

// POST registers fn on the underlying ServeMux.
func (r *GoRouterGroup) POST(relativePath string, fn HandlerFunc) {
	pattern := "/" + join(r.prefix, relativePath)
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, _ := io.ReadAll(req.Body)
		status, resp := fn(body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(resp)
	})
}
