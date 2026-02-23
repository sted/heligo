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
	get           *node
	trees         map[string]*node
	middlewares   []Middleware
	ErrorHandler  func(http.ResponseWriter, *http.Request, int, error)
	TrailingSlash bool
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
// If TrailingSlash is true, both "/path" and "/path/" will match.
func (router *Router) Handle(method string, path string, handler Handler) {
	handler = chain(handler, router.middlewares)
	router.addRoute(method, path, handler)

	if router.TrailingSlash && len(path) > 1 {
		if path[len(path)-1] == SLASH {
			router.addRoute(method, path[:len(path)-1], handler)
		} else {
			// skip paths ending with a wildcard param
			lastSlash := len(path) - 1
			for lastSlash > 0 && path[lastSlash] != SLASH {
				lastSlash--
			}
			if lastSlash < len(path)-1 && path[lastSlash+1] != STAR {
				router.addRoute(method, path+"/", handler)
			}
		}
	}
}

func (router *Router) addRoute(method string, path string, handler Handler) {
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
		if n == nil && method[0] == 'H' {
			n = router.get
		}
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

// hasPath checks if the path is registered under any method other than the given one
func (router *Router) hasPath(method string, path string) bool {
	var p params
	if method != http.MethodGet {
		if router.get != nil {
			if n := router.get.findNode(path, 0, &p); n != nil && n.handler != nil {
				return true
			}
		}
	}
	for m, tree := range router.trees {
		if m == method {
			continue
		}
		p = params{}
		if n := tree.findNode(path, 0, &p); n != nil && n.handler != nil {
			return true
		}
	}
	return false
}

// ServeHTTP complies with the standard http.Handler interface
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := Request{Request: r}
	handler := router.getHandler(r.Method, r.URL.Path, &req.params)
	if handler != nil {
		status, err := handler(r.Context(), w, req)
		if err != nil && router.ErrorHandler != nil {
			router.ErrorHandler(w, r, status, err)
		}
	} else {
		http.NotFound(w, r)
	}
}

// HasPath reports whether the given path is registered under any method
// other than the one specified. Useful for implementing 405 responses.
func (router *Router) HasPath(method string, path string) bool {
	return router.hasPath(method, path)
}

// Group creates a new sub-group of handlers, with common middlewares
func (g *Group) Group(path string, middlewares ...Middleware) *Group {
	mw := make([]Middleware, len(g.middlewares), len(g.middlewares)+len(middlewares))
	copy(mw, g.middlewares)
	return &Group{g.router, g.path + path, append(mw, middlewares...)}
}

// Handle registers a new handler under a group for method and path.
func (g *Group) Handle(method string, path string, handler Handler) {
	handler = chain(handler, g.middlewares)
	g.router.Handle(method, g.path+path, handler)
}
