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

func main() {
    router := heligo.New()
    router.Handle("GET", "/", func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
        fmt.Fprint(w, "Welcome!\n")
        return http.StatusOK, nil
    })
    router.Handle("GET", "/page/:name", func(ctx context.Context, w http.ResponseWriter, r heligo.Request) (int, error) {
       fmt.Fprintf(w, "Page %s!\n", r.Param("name"))
       return http.StatusOK, nil
    })

    http.ListenAndServe(":8080", router)
}

```