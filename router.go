package heligo

import (
	"net/http"
)

const (
	SLASH = '/'
	COLON = ':'
	STAR  = '*'
)

type Router struct {
	get          *node
	trees        map[string]*node
	middlewares  []Middleware
	ErrorHandler func(int, error)
}

type Group struct {
	router      *Router
	path        string
	middlewares []Middleware
}

// New creates a new router
func New() *Router {
	return &Router{trees: make(map[string]*node)}
}

// Use registers a global middleware.
// The middlewares are called in the order they are registered.
func (router *Router) Use(middlewares ...Middleware) {
	router.middlewares = append(router.middlewares, middlewares...)
}

// Group creates a new group of handlers, with common middlewares
func (router *Router) Group(path string, middlewares ...Middleware) *Group {
	return &Group{router, path, middlewares}
}

// Handle registers a new handler for method and path.
func (router *Router) Handle(method string, path string, handler Handler) {
	handler = chain(handler, router.middlewares)

	var n *node
	if method[0] == 'G' {
		n = router.get
		if n == nil {
			n = &node{}
			router.get = n
		}
	} else {
		n = router.trees[method]
		if n == nil {
			n = &node{}
			router.trees[method] = n
		}
	}

	var startParam int = -1
	var idxPath int
	for i := 0; i < len(path); i++ {
		if startParam == -1 {
			if path[i] == COLON || path[i] == STAR {
				startParam = i + 1
				n = n.nextNode(path[idxPath:i])
				n = n.nextNode(path[i : i+1])
			}
		} else {
			// scanning the param
			if path[i] == SLASH {
				n.param = path[startParam:i]
				startParam = -1
				idxPath = i
			}
		}
	}
	if startParam != -1 {
		n.param = path[startParam:]
	} else if idxPath <= len(path)-1 {
		n = n.nextNode(path[idxPath:])
	}
	n.handler = handler
}

func (router *Router) getHandler(method string, path string, p *params) Handler {
	var n *node
	if method[0] == 'G' {
		n = router.get
	} else {
		n = router.trees[method]
	}
	if n == nil {
		return nil
	}
	n = n.findNode(path, 0, p)
	if n != nil {
		return n.handler
	}
	return nil
}

// ServeHTTP complies with the standard http.Handler interface
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := Request{Request: r}
	handler := router.getHandler(r.Method, r.URL.Path, &req.params)
	if handler != nil {
		status, err := handler(r.Context(), w, req)
		if err != nil && router.ErrorHandler != nil {
			router.ErrorHandler(status, err)
		}
	} else {
		http.NotFound(w, r)
	}
}

// Group creates a new sub-group of handlers, with common middlewares
func (g *Group) Group(path string, middlewares ...Middleware) *Group {
	return &Group{g.router, g.path + path, append(g.middlewares, middlewares...)}
}

// Handle registers a new handler under a group for method and path.
func (g *Group) Handle(method string, path string, handler Handler) {
	handler = chain(handler, g.middlewares)
	g.router.Handle(method, g.path+path, handler)
}
