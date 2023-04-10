package httputils

import "net/http"

// Mux is an HTTP mux.
type Mux struct {
	handlers map[string]map[string]func(http.ResponseWriter, *http.Request)
}

// Handle sets a handler for a path-method pair.
func (m *Mux) Handle(path, method string, handler func(http.ResponseWriter, *http.Request)) {
	if m.handlers == nil {
		m.handlers = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	}
	methods, ok := m.handlers[path]
	if !ok {
		methods = make(map[string]func(http.ResponseWriter, *http.Request))
		m.handlers[path] = methods
	}
	methods[method] = handler
}

// ServeMux returns a stdlib ServeMux.
func (m *Mux) ServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	for path, methods := range m.handlers {
		methods := methods
		mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
			h, ok := methods[req.Method]
			if !ok {
				http.NotFound(w, req)
				return
			}
			h(w, req)
		})
	}
	return mux
}
