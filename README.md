# Heligo

Heligo is a fast and minimal HTTP request router for Go.

## Main characteristics

* Zero allocations
* Support for URL parameters (:param and *param) with precedence
* Support for middlewares and groups of handlers
* Explicit standard context in handlers
* Explicit HTTP status code and error propagation
* No internal sync.Pool usage
* No dependencies outside the standard library

## Example

```go

package main

import (
    "context"
    "fmt"
    "net/http"

    "github.com/sted/heligo"
)

func hello(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
    fmt.Fprint(w, "Welcome!\n")
    return http.StatusOK, nil
}

func page(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
    fmt.Fprintf(w, "Page %s!\n", r.Param("name"))
    return http.StatusOK, nil
}

func main() {
    router := heligo.New()
    router.Handle("GET", "/", hello)
    router.Handle("GET", "/page/:name", page)

    http.ListenAndServe(":8080", router)
}

```

## Rationale

The handler has some important differences from the standard handler:

* The context is explicit to simplify its usage (no need to extract it from the request and then reinsert it, avoiding extra allocations)
* heligo.Request wraps http.Request and gives access to the request parameters
* The HTTP status and the eventual error are returned to optimize for the fact that they are usually needed

## Groups and middlewares

You can create create groups of handlers, each with their own middlewares:

```go

projects := router.Group("/projects", DatabaseMiddleware())
projects.Handle("POST", "/", CreateProject)

```
