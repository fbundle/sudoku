package http_transport

// HandlerFunc handles a JSON-encoded request body and returns a status code and response body.
type HandlerFunc func(body []byte) (status int, resp []byte)

// Router registers POST handlers by path.
type Router interface {
	POST(path string, fn HandlerFunc)
}
